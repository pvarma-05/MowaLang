"use client";

import { useState, useEffect, useRef } from "react";
import dynamic from "next/dynamic";
import { Play, Trash2 } from "lucide-react";
import gsap from "gsap";
import type { editor } from "monaco-editor";
import type { Monaco } from "@monaco-editor/react";

const MonacoEditor = dynamic(() => import("@monaco-editor/react").then((mod) => {
  // console.log("Monaco Editor module loaded");
  return mod;
}), {
  ssr: false,
  loading: () => <div className="text-center text-gray-500 p-4">Loading editor...</div>,
});

export default function Playground() {
  const [code, setCode] = useState("");
  const [output, setOutput] = useState("");
  const [dialogue, setDialogue] = useState("");
  const [errors, setErrors] = useState<string[]>([]);
  const [isError, setIsError] = useState(false);
  const [loading, setLoading] = useState(false);
  const [wasmReady, setWasmReady] = useState(false);
  const sectionRef = useRef<HTMLElement>(null);
  const titleRef = useRef<HTMLHeadingElement>(null);
  const editorRef = useRef<HTMLDivElement>(null);
  const outputRef = useRef<HTMLDivElement>(null);
  const dialogueRef = useRef<HTMLParagraphElement>(null);
  const errorsRef = useRef<HTMLUListElement>(null);

  useEffect(() => {
    if (typeof window === "undefined") return;

    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          preloadAssets();
          loadWasm();
          observer.disconnect();
        }
      },
      { threshold: 0.1 }
    );

    if (sectionRef.current) observer.observe(sectionRef.current);

    const preloadAssets = () => {
      const assets = [
        { href: "/wasm_exec.js", as: "script" },
        { href: "/main.wasm", as: "fetch", crossOrigin: "anonymous" },
      ];
      assets.forEach(({ href, as, crossOrigin }) => {
        const link = document.createElement("link");
        link.rel = "preload";
        link.href = href;
        link.as = as;
        if (crossOrigin) link.crossOrigin = crossOrigin;
        document.head.appendChild(link);
      });
      // console.log("Preloaded WASM assets");
    };

    const loadWasm = async () => {
      if (!(window as any).Go) {
        const script = document.createElement("script");
        script.src = "/wasm_exec.js";
        script.async = true;
        script.onload = () => initializeWasm();
        script.onerror = () => {
          console.error("Failed to load wasm_exec.js");
          setOutput("Mowa, wasm_exec.js load avvaledhu!");
        };
        document.head.appendChild(script);
      } else {
        initializeWasm();
      }
    };

    const initializeWasm = async () => {
      const Go = (window as any).Go;
      if (!Go) {
        console.error("Go class not found");
        setOutput("Go class not found!");
        return;
      }

      const go = new Go();
      try {
        const response = await fetch("/main.wasm");
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        const buffer = await response.arrayBuffer();
        const module = await WebAssembly.compile(buffer);
        const instance = await WebAssembly.instantiate(module, go.importObject);
        go.run(instance);
        setWasmReady(true);
        // console.log("WASM loaded successfully");
      } catch (err) {
        console.error("WASM Load Error:", err);
        setOutput(`MowaLang Load Error: ${(err as Error).message}`);
      }
    };

    return () => observer.disconnect();
  }, []);

  useEffect(() => {
    if (!titleRef.current || !editorRef.current || !outputRef.current) return;

    const reduceMotion = window.matchMedia("(prefers-reduced-motion: reduce)").matches;

    gsap.set([titleRef.current, editorRef.current, outputRef.current], {
      autoAlpha: 0,
      y: 20,
    });

    const tl = gsap.timeline({
      defaults: { ease: "power2.out", duration: reduceMotion ? 0 : 0.5 },
    });

    tl.to(titleRef.current, { autoAlpha: 1, y: 0 })
      .to(editorRef.current, { autoAlpha: 1, y: 0 }, "-=0.2")
      .to(outputRef.current, { autoAlpha: 1, y: 0 }, "-=0.2");
  }, []);

  useEffect(() => {
    if (!dialogueRef.current && !errorsRef.current) return;

    const reduceMotion = window.matchMedia("(prefers-reduced-motion: reduce)").matches;

    gsap.set([dialogueRef.current, errorsRef.current].filter(Boolean), {
      autoAlpha: 0,
      y: 10,
    });

    gsap.to([dialogueRef.current, errorsRef.current].filter(Boolean), {
      autoAlpha: 1,
      y: 0,
      duration: reduceMotion ? 0 : 0.2,
      ease: "power2.out",
    });
  }, [dialogue, errors]);

  const handleRun = () => {
    if (!wasmReady) {
      setDialogue("");
      setErrors(["MowaLang load ayye varaku wait chey mowa!"]);
      setIsError(true);
      return;
    }

    setLoading(true);
    setOutput("");
    setDialogue("");
    setErrors([]);
    setIsError(false);

    try {
      const inputCallback = () => window.prompt("Mowa, input ivvu mowa:") || "";
      const result = (window as any).runMowa(code, inputCallback);
      const outputText = result.result || result.error || "";

      const dialogueMatch = outputText.match(/\[DIALOGUE\]((?:.|\n)*?)\[\/DIALOGUE\]/);
      let dialogueText = dialogueMatch ? dialogueMatch[1] : "";
      dialogueText = dialogueText.replace(/\\n/g, "\n");

      let remainingText = dialogueMatch ? outputText.replace(/\[DIALOGUE\]((?:.|\n)*?)\[\/DIALOGUE\]/, "").trim() : outputText;

      const errorLines = remainingText
        .split("\n")
        .filter((line: string) => line.startsWith("ln ") || line.startsWith("-:"));
      const outputLines = remainingText
        .split("\n")
        .filter((line: string) => !line.startsWith("ln ") && !line.startsWith("-:") && line !== "Mowa, errors unnai mowa:");

      setOutput(outputLines.join("\n").trim());
      setDialogue(dialogueText);
      setErrors(errorLines);
      setIsError(errorLines.length > 0 || outputText.includes("Mowa, errors unnai mowa:"));
    } catch (err: any) {
      console.error("Run Error:", err);
      setDialogue("");
      setErrors([`Error: ${err.message}`]);
      setIsError(true);
    }

    setLoading(false);
  };

  const handleClear = () => {
    setCode("");
    setOutput("");
    setDialogue("");
    setErrors([]);
    setIsError(false);
  };

  const handleEditorMount = (editor: editor.IStandaloneCodeEditor, monaco: Monaco) => {
    const textarea = editor.getDomNode()?.querySelector("textarea");
    if (textarea) {
      textarea.removeAttribute("autocomplete");
      textarea.removeAttribute("autocorrect");
      textarea.removeAttribute("autocapitalize");
      textarea.addEventListener(
        "paste",
        (e) => {
          // console.log("Paste event detected");
          e.stopPropagation();
        },
        true
      );
    }
    // Configure Monaco for mowalang language (if custom setup exists)
    monaco.languages.register({ id: "mowalang" });
    monaco.editor.setTheme("mowaTheme");
  };

  return (
    <section
      ref={sectionRef}
      className="w-full px-4 py-12 md:py-16"
      aria-labelledby="playground-title"
    >
      <h2
        ref={titleRef}
        id="playground-title"
        className="text-4xl md:text-5xl lg:text-6xl font-wg text-[#A93E39] mb-8 text-center select-none tracking-wide"
      >
        Try MOwa-Lang Online
      </h2>
      <div className="container mx-auto max-w-7xl">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <div ref={editorRef} className="bg-white rounded-xl shadow-lg overflow-hidden">
            <div className="flex justify-between items-center bg-[#A93E39] text-white p-4">
              <span className="font-outfit font-semibold select-none" aria-hidden="true">
                Code Editor
              </span>
              <div className="flex gap-3">
                <button
                  onClick={handleRun}
                  disabled={loading || !wasmReady}
                  className="px-4 py-2 bg-[#A93E39] text-white rounded-md font-outfit font-medium hover:scale-105 transition-transform duration-200 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-[#7B2D2A]"
                  aria-label={loading ? "Running code" : "Run MowaLang code"}
                  title={loading ? "Running..." : "Run code"}
                >
                  <Play size={18} aria-hidden="true" />
                  {loading ? "Running..." : "Run"}
                </button>
                <button
                  onClick={handleClear}
                  className="px-4 py-2 bg-[#A93E39] text-white rounded-md font-outfit font-medium hover:scale-105 transition-transform duration-200 flex items-center gap-2 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-[#7B2D2A]"
                  aria-label="Clear editor and output"
                  title="Clear editor"
                >
                  <Trash2 size={18} aria-hidden="true" />
                  Clear
                </button>
              </div>
            </div>
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
                minimap: { enabled: false },
                wordWrap: "on",
                automaticLayout: true,
                tabSize: 2,
                padding: { top: 10, bottom: 10 },
                scrollBeyondLastLine: false,
                lineNumbers: "on",
                renderLineHighlight: "line",
                scrollbar: { vertical: "auto", horizontal: "auto" },
                cursorBlinking: "smooth",
                accessibilitySupport: "on",
                domReadOnly: false,
                readOnly: false,
              }}
              aria-label="MowaLang code editor with paste support"
            />
          </div>
          <div ref={outputRef} className="bg-white rounded-xl shadow-lg overflow-hidden">
            <div className="bg-[#A93E39] text-white p-4 font-outfit font-semibold select-none" aria-hidden="true">
              Output
            </div>
            <div
              className="p-4 font-outfit text-md text-gray-800 whitespace-pre-wrap min-h-[450px] max-h-[450px] overflow-auto"
              role="log"
              aria-live="polite"
              aria-label="MowaLang output"
            >
              {output && (
                <pre className="mb-2" aria-label="Program output">
                  {output}
                </pre>
              )}
              {dialogue && (
                <p
                  ref={dialogueRef}
                  className={`font-outfit text-md italic bg-gray-200 rounded-md px-3 py-2 mb-2 !whitespace-pre-wrap ${
                    isError ? "text-red-600" : "text-green-600"
                  }`}
                  aria-label={isError ? "Error dialogue" : "Success dialogue"}
                  aria-multiline="true"
                >
                  {dialogue}
                </p>
              )}
              {errors.length > 0 && (
                <ul
                  ref={errorsRef}
                  className="font-outfit text-md text-red-600"
                  aria-label="Error messages"
                >
                  {errors.map((err, index) => (
                    <li key={index} className="mb-1">
                      {err}
                    </li>
                  ))}
                </ul>
              )}
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
