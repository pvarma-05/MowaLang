"use client";

import { useState, useEffect } from "react";
import dynamic from "next/dynamic";
import { motion } from "framer-motion";
import { Play, Trash2 } from "lucide-react";
import Link from "next/link";

const MonacoEditor = dynamic(() => import("@monaco-editor/react"), {
  ssr: false,
  loading: () => <div className="text-center text-gray-500">Loading editor...</div>,
});

export default function Home() {
  const [code, setCode] = useState("");
  const [output, setOutput] = useState("Loading WebAssembly module...");
  const [loading, setLoading] = useState(false);
  const [wasmReady, setWasmReady] = useState(false);

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
      if (!Go) return setOutput("Go class not found even after wasm_exec.js load!");

      const go = new Go();
      WebAssembly.instantiateStreaming(fetch("/mowalang.wasm"), go.importObject)
        .then((result) => {
          go.run(result.instance);
          setWasmReady(true);
          setOutput((prev) => (prev === "Loading WebAssembly module..." ? "Output appears here..." : prev));
        })
        .catch((err) => setOutput(`MowaLang Load Error: ${err.message}`));
    };

    loadWasm();
  }, []);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!wasmReady) return setOutput("MowaLang load ayye varaku wait chey mowa!");

    setLoading(true);
    setOutput("");

    try {
      const inputCallback = () => window.prompt("Mowa, input ivvu mowa:") || "";
      const result = (window as any).runMowa(code, inputCallback);
      setOutput(result.result || result.error);
    } catch (err: any) {
      setOutput(`Error: ${err.message}`);
    }

    setLoading(false);
  }

  const handleClear = () => {
    setCode("");
    setOutput("Output appears here...");
  };

  return (
    <main className="flex flex-col min-h-screen bg-gradient-to-b from-[#FDF2EC] to-[#FFF7F3] text-black">
      <header className="w-full bg-[#A93E39] text-white py-4 shadow-md min-h-[60px]">
        <div className="container mx-auto px-4 flex justify-between items-center">
          <Link href="/" className="text-4xl font-wg select-none">MOwa-Lang</Link>
          <nav className="flex gap-4 select-none">
            <a href="/docs" className="hover:underline font-outfit">Docs</a>
            <a href="https://github.com/pvarma-05/MowaLang" target="_blank" rel="noopener noreferrer" className="hover:underline font-outfit">GitHub</a>
          </nav>
        </div>
      </header>

      <section className="container mx-auto px-4 py-8 flex-grow">
        <div className="text-center mb-8">
          <motion.h2 initial={{ opacity: 0, y: -20 }} animate={{ opacity: 1, y: 0 }} transition={{ duration: 0.5 }} className="text-4xl md:text-6xl font-outfit text-[#A93E39] select-none">
            <span className="font-wg">MOwa-Lang PlaygrOUnd</span>
          </motion.h2>
          <p className="mt-2 text-lg font-outfit text-gray-700">Code in Telugu with dramatic flair!</p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <div className="flex flex-col">
            <div className="bg-white rounded-lg shadow-lg overflow-hidden">
              <div className="flex justify-between items-center bg-[#A93E39] text-white p-3">
                <span className="font-outfit font-semibold select-none">Code Editor</span>
                <div className="flex gap-2 select-none">
                  <button onClick={handleSubmit} disabled={loading || !wasmReady} className="cursor-pointer flex items-center gap-2 px-3 py-1 bg-white text-[#A93E39] rounded-md font-outfit font-semibold hover:bg-gray-100 disabled:opacity-50">
                    <Play size={16} /> {loading ? "Running..." : "Run"}
                  </button>
                  <button onClick={handleClear} className="cursor-pointer flex items-center gap-2 px-3 py-1 bg-white text-[#A93E39] rounded-md font-outfit font-semibold hover:bg-gray-100">
                    <Trash2 size={16} /> Clear
                  </button>
                </div>
              </div>
              <MonacoEditor
                height="500px"
                language="mowalang"
                theme="mowaTheme"
                value={code}
                onChange={(value) => setCode(value || "")}
                options={{
                  fontSize: 16,
                  fontFamily: "'Outfit', monospace",
                  cursorStyle: "line",
                  minimap: { enabled: false },
                  wordWrap: "on",
                  automaticLayout: true,
                  tabSize: 2,
                  padding: { top: 10, bottom: 10 },
                  scrollBeyondLastLine: false,
                  lineNumbers: "on",
                  renderLineHighlight: "all",
                  scrollbar: { vertical: "visible", horizontal: "visible" },
                }}
                beforeMount={(monaco) => {
                  monaco.languages.register({ id: "mowalang" });
                  monaco.languages.setMonarchTokensProvider("mowalang", {
                    keywords: ["idhi", "theesko", "mowa", "okavela", "ayithe", "ledha", "ledhante", "enchuko", "case", "default", "pani", "ichey", "nijam", "abadham", "aagipo", "kaani", "varaku", "rakam", "number", "string", "boolean"],
                    operators: ["+", "-", "*", "/", "%", "**", "=", "==", "!=", "<", "<=", ">", ">=", "&&", "||", "++", "--", "+=", "-=", "*=", "/=", "%="],
                    tokenizer: {
                      root: [
                        [/[a-zA-Z_][a-zA-Z0-9_]*/, { cases: { "@keywords": "keyword.mowalang", "@default": "identifier.mowalang" } }],
                        [/[+\-*/%=&|<>!]=?|\*\*|[(){}\[\],:]/, "operator.mowalang"],
                        [/"[^"]*"|'[^']*'/, "string.mowalang"],
                        [/\d+(\.\d+)?/, "number.mowalang"],
                        [/\/\/.*/, "comment.mowalang"],
                        [/\s+/, "white"],
                      ],
                    },
                  });
                  monaco.editor.defineTheme("mowaTheme", {
                    base: "vs",
                    inherit: true,
                    rules: [
                      { token: "keyword.mowalang", foreground: "#D73A49" },
                      { token: "string.mowalang", foreground: "#22863A" },
                      { token: "number.mowalang", foreground: "#005CC5" },
                      { token: "comment.mowalang", foreground: "#6A737D" },
                      { token: "identifier.mowalang", foreground: "#24292E" },
                      { token: "operator.mowalang", foreground: "#D73A49" },
                    ],
                    colors: {
                      "editor.background": "#FAFAFA",
                      "editor.foreground": "#24292E",
                      "editorLineNumber.foreground": "#D73A49",
                      "editorLineNumber.activeForeground": "#B9252E",
                      "editor.selectionBackground": "#D73A4933",
                      "editorCursor.foreground": "#D73A49",
                    },
                  });
                }}
              />
            </div>
            {!wasmReady && <p className="mt-2 text-[#A93E39] font-outfit">Loading MowaLang WASM module...</p>}
          </div>

          <div className="flex flex-col">
            <div className="bg-[#FDF2EC] rounded-lg shadow-lg overflow-hidden">
              <div className="bg-[#A93E39] text-white p-4 font-outfit font-semibold select-none">Output</div>
              <pre className="p-4 font-outfit text-md h-[500px] overflow-auto text-gray-800 whitespace-pre-wrap break-words">{output}</pre>
            </div>
          </div>
        </div>
      </section>

      <footer className="w-full bg-[#A93E39] text-white py-4 mt-auto min-h-[60px] select-none">
        <div className="container mx-auto px-4 text-center font-outfit">
          <p>Built with 🤍 by <a href="https://github.com/pvarma-05" target="_blank" rel="noopener noreferrer" className="underline">Pradeep Varma</a></p>
          <p className="mt-1">MowaLang © {new Date().getFullYear()}</p>
        </div>
      </footer>
    </main>
  );
}
