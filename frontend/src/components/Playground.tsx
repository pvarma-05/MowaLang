"use client";

import { useState, useEffect, useRef } from "react";
import dynamic from "next/dynamic";
import { Play, Trash2 } from "lucide-react";
import gsap from "gsap";

const MonacoEditor = dynamic(() => import("@monaco-editor/react"), {
  ssr: false,
  loading: () => <div className="text-center text-gray-500">Loading editor...</div>,
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

    const loadWasm = async () => {
      if (!(window as any).Go) {
        const script = document.createElement("script");
        script.src = "/wasm_exec.js";
        script.async = true;
        script.onload = () => initializeWasm();
        script.onerror = () => setOutput("Mowa, wasm_exec.js load avvaledhu!");
        document.head.appendChild(script);
      } else {
        initializeWasm();
      }
    };

    const initializeWasm = () => {
      const Go = (window as any).Go;
      if (!Go) return setOutput("Go class not found!");

      const go = new Go();
      WebAssembly.instantiateStreaming(fetch("/mowalang.wasm"), go.importObject)
        .then((result) => {
          go.run(result.instance);
          setWasmReady(true);
        })
        .catch((err) => setOutput(`MowaLang Load Error: ${err.message}`));
    };

    loadWasm();
  }, []);

  useEffect(() => {
    if (!titleRef.current || !editorRef.current || !outputRef.current) return;

    const reduceMotion = window.matchMedia("(prefers-reduced-motion: reduce)").matches;

    gsap.set([titleRef.current, editorRef.current, outputRef.current], {
      autoAlpha: 0,
      y: 50,
    });

    const tl = gsap.timeline({
      defaults: { ease: "power2.out", duration: reduceMotion ? 0 : 0.8 },
    });

    tl.to(titleRef.current, {
      autoAlpha: 1,
      y: 0,
      duration: reduceMotion ? 0 : 1,
    })
      .to(
        editorRef.current,
        {
          autoAlpha: 1,
          y: 0,
          stagger: 0.2,
        },
        "-=0.4"
      )
      .to(
        outputRef.current,
        {
          autoAlpha: 1,
          y: 0,
          stagger: 0.2,
        },
        "-=0.4"
      );
  }, []);

  useEffect(() => {
    if (!dialogueRef.current && !errorsRef.current) return;

    const reduceMotion = window.matchMedia("(prefers-reduced-motion: reduce)").matches;

    gsap.set([dialogueRef.current, errorsRef.current].filter(Boolean), {
      autoAlpha: 0,
      y: 20,
    });

    gsap.to([dialogueRef.current, errorsRef.current].filter(Boolean), {
      autoAlpha: 1,
      y: 0,
      duration: reduceMotion ? 0 : 0.5,
      ease: "power2.out",
      stagger: 0.1,
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
      dialogueText = dialogueText.replace(/\\n/g, '\n');

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

  return (
    <section
      ref={sectionRef}
      className="w-full px-4 py-12  md:py-16"
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
          {/* Editor Panel */}
          <div ref={editorRef} className="bg-white rounded-xl shadow-lg overflow-hidden">
            <div className="flex justify-between items-center bg-[#A93E39] text-white p-4">
              <span className="font-outfit font-semibold select-none" aria-hidden="true">
                Code Editor
              </span>
              <div className="flex gap-3">
                <button
                  onClick={handleRun}
                  disabled={loading || !wasmReady}
                  className="px-4 py-2 bg-[#A93E39] text-white rounded-md font-outfit font-medium hover:scale-105 transition-transform duration-300 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-[#7B2D2A] cursor-pointer"
                  aria-label={loading ? "Running code" : "Run MowaLang code"}
                  title={loading ? "Running..." : "Run code"}
                >
                  <Play size={18} aria-hidden="true" />
                  {loading ? "Running..." : "Run"}
                </button>
                <button
                  onClick={handleClear}
                  className="cursor-pointer px-4 py-2 bg-[#A93E39] text-white rounded-md font-outfit font-medium hover:scale-105 transition-transform duration-300 flex items-center gap-2 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-[#7B2D2A]"
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
              options={{
                fontSize: 16,
                fontFamily: "'Outfit', monospace",
                fontLigatures: false,
                minimap: { enabled: false },
                wordWrap: "on",
                automaticLayout: true,
                tabSize: 2,
                padding: { top: 10, bottom: 10 },
                scrollBeyondLastLine: false,
                lineNumbers: "on",
                renderLineHighlight: "line",
                scrollbar: { vertical: "auto", horizontal: "auto" },
                smoothScrolling: true,
                cursorBlinking: "smooth",
                accessibilitySupport: "on",
              }}
              aria-label="MowaLang code editor"
            />
          </div>

          {/* Output Panel */}
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
