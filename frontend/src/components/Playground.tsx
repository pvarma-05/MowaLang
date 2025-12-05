"use client";
import {
  useState,
  useEffect,
  useRef,
  useImperativeHandle,
  forwardRef,
  useCallback,
  useMemo,
} from "react";
import dynamic from "next/dynamic";
import { Play, Trash2, Loader2, CheckCircle2, AlertCircle } from "lucide-react";
import gsap from "gsap";
import type { editor } from "monaco-editor";
import type { Monaco } from "@monaco-editor/react";

const MonacoEditor = dynamic(
  () =>
    import("@monaco-editor/react").then((mod) => {
      return mod;
    }),
  {
    ssr: false,
    loading: () => (
      <div className="flex items-center justify-center h-[450px] bg-gray-50">
        <div className="flex flex-col items-center gap-3">
          <Loader2 className="w-8 h-8 text-[#A93E39] animate-spin" />
          <p className="text-gray-600 font-outfit">Loading editor...</p>
        </div>
      </div>
    ),
  }
);

export interface PlaygroundHandle {
  setCode: (code: string) => void;
}

type ExecutionStatus = "idle" | "loading" | "running" | "success" | "error";

const Playground = forwardRef<PlaygroundHandle>((props, ref) => {
const [code, setCode] = useState(`pani greet(name : string):string{
  ichey "Hello "+name+"!,Thanks For Using MowaLang."
}
idhi name : string;
theesko name;
mowa greet(name);`);
  const [output, setOutput] = useState("");
  const [dialogue, setDialogue] = useState("");
  const [errors, setErrors] = useState<string[]>([]);
  const [isError, setIsError] = useState(false);
  const [executionStatus, setExecutionStatus] = useState<ExecutionStatus>("idle");
  const [wasmReady, setWasmReady] = useState(false);
  const [wasmLoadProgress, setWasmLoadProgress] = useState(0);

  const sectionRef = useRef<HTMLElement>(null);
  const titleRef = useRef<HTMLHeadingElement>(null);
  const editorRef = useRef<HTMLDivElement>(null);
  const outputRef = useRef<HTMLDivElement>(null);
  const monacoEditorRef = useRef<editor.IStandaloneCodeEditor | null>(null);
  const wasmInitializedRef = useRef(false);

  useImperativeHandle(ref, () => ({
    setCode: (newCode: string) => {
      setCode(newCode);
      if (monacoEditorRef.current) {
        monacoEditorRef.current.setValue(newCode);
        monacoEditorRef.current.focus();
      }
    },
  }));

  // Memoize status styling
  const statusIndicator = useMemo(() => {
    const statusConfig = {
      idle: { icon: null, color: "text-gray-400", bg: "bg-gray-100" },
      loading: { icon: Loader2, color: "text-blue-500", bg: "bg-blue-50" },
      running: { icon: Loader2, color: "text-yellow-500", bg: "bg-yellow-50" },
      success: { icon: CheckCircle2, color: "text-green-500", bg: "bg-green-50" },
      error: { icon: AlertCircle, color: "text-red-500", bg: "bg-red-50" },
    };
    return statusConfig[executionStatus];
  }, [executionStatus]);

  // WASM Loading with progress
  useEffect(() => {
    if (typeof window === "undefined" || wasmInitializedRef.current) return;

    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          initializeWasmSystem();
          observer.disconnect();
        }
      },
      { threshold: 0.1, rootMargin: "50px" }
    );

    if (sectionRef.current) observer.observe(sectionRef.current);

    const initializeWasmSystem = async () => {
      setExecutionStatus("loading");
      setWasmLoadProgress(10);

      try {
        // Preload assets
        const preloadPromises = [
          new Promise((resolve) => {
            const script = document.createElement("link");
            script.rel = "preload";
            script.href = "/wasm_exec.js";
            script.as = "script";
            document.head.appendChild(script);
            resolve(true);
          }),
          new Promise((resolve) => {
            const link = document.createElement("link");
            link.rel = "preload";
            link.href = "/main.wasm";
            link.as = "fetch";
            link.crossOrigin = "anonymous";
            document.head.appendChild(link);
            resolve(true);
          }),
        ];

        await Promise.all(preloadPromises);
        setWasmLoadProgress(30);

        // Load wasm_exec.js
        if (!(window as any).Go) {
          await new Promise<void>((resolve, reject) => {
            const script = document.createElement("script");
            script.src = "/wasm_exec.js";
            script.async = true;
            script.onload = () => {
              setWasmLoadProgress(60);
              resolve();
            };
            script.onerror = () => reject(new Error("Failed to load wasm_exec.js"));
            document.head.appendChild(script);
          });
        } else {
          setWasmLoadProgress(60);
        }

        // Initialize WASM
        const Go = (window as any).Go;
        if (!Go) throw new Error("Go class not found");

        const go = new Go();
        setWasmLoadProgress(70);

        const response = await fetch("/main.wasm");
        if (!response.ok) throw new Error(`HTTP ${response.status}`);

        const buffer = await response.arrayBuffer();
        setWasmLoadProgress(85);

        const module = await WebAssembly.compile(buffer);
        const instance = await WebAssembly.instantiate(module, go.importObject);
        setWasmLoadProgress(95);

        go.run(instance);
        setWasmLoadProgress(100);

        wasmInitializedRef.current = true;
        setWasmReady(true);
        setExecutionStatus("idle");
      } catch (err) {
        console.error("WASM Load Error:", err);
        setOutput(`MowaLang Load Error: ${(err as Error).message}`);
        setExecutionStatus("error");
      }
    };

    return () => observer.disconnect();
  }, []);

  // Smooth entrance animations
  useEffect(() => {
    if (!titleRef.current || !editorRef.current || !outputRef.current) return;

    const reduceMotion = window.matchMedia("(prefers-reduced-motion: reduce)").matches;

    gsap.set([titleRef.current, editorRef.current, outputRef.current], {
      autoAlpha: 0,
      y: 30,
    });

    const tl = gsap.timeline({
      defaults: { ease: "power3.out", duration: reduceMotion ? 0 : 0.8 },
    });

    tl.to(titleRef.current, { autoAlpha: 1, y: 0 })
      .to(editorRef.current, { autoAlpha: 1, y: 0 }, "-=0.5")
      .to(outputRef.current, { autoAlpha: 1, y: 0 }, "-=0.5");
  }, []);

  // Run code with smooth feedback
  const handleRun = useCallback(() => {
    if (!wasmReady) {
      setDialogue("");
      setErrors(["MowaLang load ayye varaku wait chey mowa!"]);
      setIsError(true);
      setExecutionStatus("error");
      return;
    }

    setExecutionStatus("running");
    setOutput("");
    setDialogue("");
    setErrors([]);
    setIsError(false);

    // Add slight delay for smooth transition
    setTimeout(() => {
      try {
        const inputCallback = () => window.prompt("Mowa, input ivvu mowa:") || "";
        const result = (window as any).runMowa(code, inputCallback);
        const outputText = result.result || result.error || "";

        const dialogueMatch = outputText.match(/\[DIALOGUE\]((?:.|\n)*?)\[\/DIALOGUE\]/);
        let dialogueText = dialogueMatch ? dialogueMatch[1] : "";
        dialogueText = dialogueText.replace(/\\n/g, "\n");

        let remainingText = dialogueMatch
          ? outputText.replace(/\[DIALOGUE\]((?:.|\n)*?)\[\/DIALOGUE\]/, "").trim()
          : outputText;

        const errorLines = remainingText
          .split("\n")
          .filter((line: string) => line.startsWith("ln ") || line.startsWith("-:"));

        const outputLines = remainingText
          .split("\n")
          .filter(
            (line: string) =>
              !line.startsWith("ln ") &&
              !line.startsWith("-:") &&
              line !== "Mowa, errors unnai mowa:"
          );

        setOutput(outputLines.join("\n").trim());
        setDialogue(dialogueText);
        setErrors(errorLines);

        const hasErrors = errorLines.length > 0 || outputText.includes("Mowa, errors unnai mowa:");
        setIsError(hasErrors);
        setExecutionStatus(hasErrors ? "error" : "success");
      } catch (err: any) {
        console.error("Run Error:", err);
        setDialogue("");
        setErrors([`Error: ${err.message}`]);
        setIsError(true);
        setExecutionStatus("error");
      }
    }, 100);
  }, [wasmReady, code]);

  // Clear with animation
  const handleClear = useCallback(() => {
    const reduceMotion = window.matchMedia("(prefers-reduced-motion: reduce)").matches;

    if (!reduceMotion && outputRef.current) {
      gsap.to(outputRef.current.children, {
        autoAlpha: 0,
        y: -10,
        duration: 0.2,
        stagger: 0.05,
        onComplete: () => {
          setCode("");
          setOutput("");
          setDialogue("");
          setErrors([]);
          setIsError(false);
          setExecutionStatus("idle");

          gsap.to(outputRef.current!.children, {
            autoAlpha: 1,
            y: 0,
            duration: 0.3,
          });
        },
      });
    } else {
      setCode("");
      setOutput("");
      setDialogue("");
      setErrors([]);
      setIsError(false);
      setExecutionStatus("idle");
    }
  }, []);

  const handleEditorMount = useCallback((editor: editor.IStandaloneCodeEditor, monaco: Monaco) => {
    monacoEditorRef.current = editor;

    const textarea = editor.getDomNode()?.querySelector("textarea");
    if (textarea) {
      textarea.removeAttribute("autocomplete");
      textarea.removeAttribute("autocorrect");
      textarea.removeAttribute("autocapitalize");
      textarea.setAttribute("spellcheck", "false");
    }

    monaco.languages.register({ id: "mowalang" });

    // Custom theme
    monaco.editor.defineTheme("mowaTheme", {
      base: "vs",
      inherit: true,
      rules: [
        { token: "keyword", foreground: "A93E39", fontStyle: "bold" },
        { token: "string", foreground: "2E7D32" },
        { token: "number", foreground: "1565C0" },
        { token: "comment", foreground: "757575", fontStyle: "italic" },
      ],
      colors: {
        "editor.background": "#FAFAFA",
        "editor.lineHighlightBackground": "#F5F5F5",
        "editorCursor.foreground": "#A93E39",
        "editor.selectionBackground": "#FFE0DE",
      },
    });

    monaco.editor.setTheme("mowaTheme");
  }, []);

  const StatusIcon = statusIndicator.icon;

  return (
    <section
      ref={sectionRef}
      className="w-full px-4 py-12 md:py-16 relative"
      aria-labelledby="playground-title"
    >
      {/* Background decoration */}
      <div className="absolute inset-0 overflow-hidden pointer-events-none">
        <div className="absolute top-0 right-1/4 w-72 h-72 bg-[#A93E39] opacity-5 rounded-full blur-3xl"></div>
        <div className="absolute bottom-0 left-1/4 w-96 h-96 bg-[#A93E39] opacity-5 rounded-full blur-3xl"></div>
      </div>

      <div className="relative z-10">
        <h2
          ref={titleRef}
          id="playground-title"
          className="text-4xl md:text-5xl lg:text-6xl font-wg text-[#A93E39] mb-8 text-center select-none tracking-wide drop-shadow-sm"
        >
          Try MOwa-Lang Online
        </h2>

        {/* WASM Loading Progress */}
        {!wasmReady && executionStatus === "loading" && (
          <div className="container mx-auto max-w-7xl mb-6">
            <div className="bg-white rounded-lg shadow-md p-6">
              <div className="flex items-center gap-3 mb-3">
                <Loader2 className="w-5 h-5 text-[#A93E39] animate-spin" />
                <span className="font-outfit text-gray-700">
                  Loading MowaLang... {wasmLoadProgress}%
                </span>
              </div>
              <div className="w-full bg-gray-200 rounded-full h-2 overflow-hidden">
                <div
                  className="bg-[#A93E39] h-full transition-all duration-300 ease-out"
                  style={{ width: `${wasmLoadProgress}%` }}
                ></div>
              </div>
            </div>
          </div>
        )}

        <div className="container mx-auto max-w-7xl">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            {/* Editor Panel */}
            <div
              ref={editorRef}
              className="bg-white rounded-xl shadow-xl overflow-hidden transition-all duration-300 hover:shadow-2xl border border-gray-100"
            >
              <div className="flex justify-between items-center bg-gradient-to-r from-[#A93E39] to-[#8B3230] text-white p-4">
                <div className="flex items-center gap-3">
                  <span className="font-outfit font-semibold select-none">Code Editor</span>
                  {StatusIcon && (
                    <div className="flex items-center gap-1.5">
                      <StatusIcon className={`w-4 h-4 ${statusIndicator.color} ${executionStatus === 'running' || executionStatus === 'loading' ? 'animate-spin' : ''}`} />
                    </div>
                  )}
                </div>
                <div className="flex gap-3">
                  <button
                    onClick={handleRun}
                    disabled={executionStatus === "running" || !wasmReady}
                    className="px-4 py-2 bg-white/20 backdrop-blur-sm text-white rounded-lg font-outfit font-medium hover:bg-white/30 transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2 focus:outline-none focus:ring-2 focus:ring-white/50 hover:scale-105 active:scale-95"
                    aria-label={executionStatus === "running" ? "Running code" : "Run MowaLang code"}
                    title={executionStatus === "running" ? "Running..." : "Run code"}
                  >
                    {executionStatus === "running" ? (
                      <Loader2 size={18} className="animate-spin" aria-hidden="true" />
                    ) : (
                      <Play size={18} aria-hidden="true" />
                    )}
                    {executionStatus === "running" ? "Running..." : "Run"}
                  </button>
                  <button
                    onClick={handleClear}
                    className="px-4 py-2 bg-white/20 backdrop-blur-sm text-white rounded-lg font-outfit font-medium hover:bg-white/30 transition-all duration-200 flex items-center gap-2 focus:outline-none focus:ring-2 focus:ring-white/50 hover:scale-105 active:scale-95"
                    aria-label="Clear editor and output"
                    title="Clear editor"
                  >
                    <Trash2 size={18} aria-hidden="true" />
                    Clear
                  </button>
                </div>
              </div>
              <div className="relative">
                <MonacoEditor
                  height="450px"
                  language="mowalang"
                  theme="mowaTheme"
                  value={code}
                  onChange={(value) => setCode(value || "")}
                  onMount={handleEditorMount}
                  options={{
                    fontSize: 16,
                    fontFamily: "'Outfit', monospace",
                    fontLigatures: true,
                    minimap: { enabled: false },
                    wordWrap: "on",
                    automaticLayout: true,
                    tabSize: 2,
                    padding: { top: 16, bottom: 16 },
                    scrollBeyondLastLine: false,
                    lineNumbers: "on",
                    renderLineHighlight: "all",
                    scrollbar: {
                      vertical: "auto",
                      horizontal: "auto",
                      verticalScrollbarSize: 10,
                      horizontalScrollbarSize: 10,
                    },
                    cursorBlinking: "smooth",
                    cursorSmoothCaretAnimation: "on",
                    smoothScrolling: true,
                    accessibilitySupport: "on",
                    suggest: { showWords: false },
                    quickSuggestions: false,
                  }}
                  aria-label="MowaLang code editor"
                />
              </div>
            </div>

            {/* Output Panel */}
            <div
              ref={outputRef}
              className="bg-white rounded-xl shadow-xl overflow-hidden transition-all duration-300 hover:shadow-2xl border border-gray-100"
            >
              <div className="bg-gradient-to-r from-[#A93E39] to-[#8B3230] text-white p-4 font-outfit font-semibold select-none flex items-center justify-between">
                <span>Output</span>
                {executionStatus !== "idle" && (
                  <span className={`text-xs px-3 py-1 rounded-full ${statusIndicator.bg} ${statusIndicator.color} font-medium backdrop-blur-sm`}>
                    {executionStatus === "loading" && "Loading WASM"}
                    {executionStatus === "running" && "Running..."}
                    {executionStatus === "success" && "Success"}
                    {executionStatus === "error" && "Error"}
                  </span>
                )}
              </div>
              <div
                className="p-6 font-outfit text-md text-gray-800 whitespace-pre-wrap min-h-[450px] max-h-[450px] overflow-auto bg-gradient-to-br from-gray-50 to-white scrollbar-thin scrollbar-thumb-[#A93E39] scrollbar-track-gray-200"
                role="log"
                aria-live="polite"
                aria-label="MowaLang output"
              >
                {!output && !dialogue && errors.length === 0 && executionStatus !== "loading" && (
                  <div className="flex flex-col items-center justify-center h-full text-gray-400">
                    <Play className="w-12 h-12 mb-3 opacity-50" />
                    <p className="text-sm">Run your code to see output</p>
                  </div>
                )}

                {output && (
                  <pre className="mb-3 p-3 bg-white rounded-lg shadow-sm border border-gray-100 animate-fadeIn">
                    {output}
                  </pre>
                )}

                {dialogue && (
                  <div
                    className={`font-outfit text-md p-4 rounded-lg shadow-sm mb-3 !whitespace-pre-wrap animate-fadeIn ${
                      isError
                        ? "bg-red-50 text-red-700 border border-red-200"
                        : "bg-green-50 text-green-700 border border-green-200"
                    }`}
                  >
                    <div className="flex items-start gap-2">
                      {isError ? (
                        <AlertCircle className="w-5 h-5 mt-0.5 flex-shrink-0" />
                      ) : (
                        <CheckCircle2 className="w-5 h-5 mt-0.5 flex-shrink-0" />
                      )}
                      <p className="italic flex-1">{dialogue}</p>
                    </div>
                  </div>
                )}

                {errors.length > 0 && (
                  <div className="bg-red-50 border border-red-200 rounded-lg p-4 animate-fadeIn">
                    <div className="flex items-start gap-2 mb-2">
                      <AlertCircle className="w-5 h-5 text-red-600 mt-0.5 flex-shrink-0" />
                      <h4 className="font-semibold text-red-700">Errors:</h4>
                    </div>
                    <ul className="font-outfit text-sm text-red-600 space-y-1 ml-7">
                      {errors.map((err, index) => (
                        <li key={index} className="font-mono">
                          {err}
                        </li>
                      ))}
                    </ul>
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
});

Playground.displayName = "Playground";
export default Playground;
