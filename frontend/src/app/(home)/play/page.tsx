"use client";

import { useState, useEffect } from "react";
import dynamic from "next/dynamic";

const MonacoEditor = dynamic(() => import("@monaco-editor/react"), {
  ssr: false,
});

export default function Home() {
  const [code, setCode] = useState("");
  const [output, setOutput] = useState("Loading WebAssembly module...");
  const [loading, setLoading] = useState(false);
  const [wasmReady, setWasmReady] = useState(false);

  // Define MowaLang theme
  const defineMowaTheme = () => {
    if (typeof window === "undefined") return;

    import('monaco-editor').then((monaco) => {
      monaco.editor.defineTheme('mowaTheme', {
        base: 'vs',
        inherit: true,
        rules: [
          { token: 'keyword', foreground: '#A93E39', fontStyle: 'bold' },
          { token: 'string', foreground: '#006400' },
          { token: 'number', foreground: '#0000FF' },
          { token: 'comment', foreground: '#6A9955' },
          { token: 'identifier', foreground: '#000000' },
        ],
        colors: {
          'editor.background': '#FDF2EC',
          'editor.foreground': '#000000',
          'editorLineNumber.foreground': '#A93E39',
          'editorLineNumber.activeForeground': '#7B2D2A',
          'editor.selectionBackground': '#A93E3933',
          'editorCursor.foreground': '#A93E39',
        },
      });
    });
  };

  // Load Monaco and WASM setup only on client side
  useEffect(() => {
    if (typeof window === "undefined") return;

    // Define theme
    defineMowaTheme();

    // Define MowaLang syntax highlighting
    import('monaco-editor').then((monaco) => {
      monaco.languages.register({ id: 'mowalang' });
      monaco.languages.setMonarchTokensProvider('mowalang', {
        tokenizer: {
          root: [
            [/(idhi|theesko|mowa|okavela|ayithe|ledha|ledhante|enchuko|case|default|pani|ichey|nijam|abadham|aagipo|kaani|varaku|rakam)/, 'keyword'],
            [/"[^"]*"|'[^']*'/, 'string'],
            [/\d+(\.\d+)?/, 'number'],
            [/[a-zA-Z_][a-zA-Z0-9_]*/, 'identifier'],
            [/\/\/.*/, 'comment'],
          ],
        },
      });
    });

    // Load wasm_exec.js and initialize WASM
    const loadWasm = async () => {
      if (!(window as any).Go) {
        const script = document.createElement("script");
        script.src = "/wasm_exec.js";
        script.async = true;
        script.onload = () => {
          initializeWasm();
        };
        script.onerror = () => {
          setOutput("Mowa, wasm_exec.js load avvaledhu!");
        };
        document.head.appendChild(script);
      } else {
        initializeWasm();
      }
    };

    const initializeWasm = () => {
      const Go = (window as any).Go;
      if (!Go) {
        setOutput("Mowa, Go class dhorakaledhu, wasm_exec.js load ayina tharvatha kooda!");
        return;
      }

      const go = new Go();
      if (!go.importObject) {
        setOutput("Mowa, Go.importObject anedhi asala ledhu mowa!");
        return;
      }

      WebAssembly.instantiateStreaming(fetch("/mowalang.wasm"), go.importObject)
        .then((result) => {
          go.run(result.instance);
          setWasmReady(true);
          setOutput("Output appears here...");
        })
        .catch((err) => {
          setOutput(`MowaLang Load Error: ${err.message}`);
        });
    };

    loadWasm();
  }, []);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!wasmReady) {
      setOutput("MowaLang load ayye varaku wait chey mowa!");
      return;
    }
    setLoading(true);
    setOutput("");

    try {
      const inputCallback = (prompt: string) => {
        return window.prompt("Mowa, input ivvu mowa:") || "";
      };
      const result = (window as any).runMowa(code, inputCallback);
      setOutput(result.result || result.error);
    } catch (err: any) {
      setOutput(`Error: ${err.message}`);
    }

    setLoading(false);
  }

  return (
    <main className="flex flex-col items-center min-h-screen text-black px-4 bg-[#FDF2EC] py-8">
      <h1 className="text-5xl md:text-7xl my-8 font-wg text-[#A93E39] tracking-wider">
        MOwa-Lang PLAYGROUND
      </h1>
      <form
        onSubmit={handleSubmit}
        className="w-full max-w-7xl flex flex-col lg:flex-row gap-4"
      >
        <div className="w-full lg:w-1/2">
          <MonacoEditor
            height="500px"
            language="mowalang"
            theme="mowaTheme"
            value={code}
            onChange={(value) => setCode(value || "")}
            options={{
              fontSize: 16,
              fontFamily: "'Outfit', sans-serif",
              minimap: { enabled: false },
              wordWrap: "on",
              automaticLayout: true,
              tabSize: 2,
              padding: { top: 10, bottom: 10 },
              scrollBeyondLastLine: false,
            }}
          />
          <div className="mt-2">
            <button
              type="submit"
              disabled={loading || !wasmReady}
              className="cursor-pointer px-6 py-2 bg-[#A93E39] text-white font-semibold font-outfit rounded-lg hover:scale-105 transition-all duration-300 disabled:bg-[#7B2D2A] disabled:cursor-not-allowed"
            >
              {loading ? "Running..." : "Run Code"}
            </button>
            {!wasmReady && (
              <span className="ml-4 text-[#A93E39] font-outfit">
                Loading MowaLang...
              </span>
            )}
          </div>
        </div>
        <pre
          className="w-full lg:w-1/2 p-6 bg-[#FFF7F3] text-black font-outfit rounded-lg shadow-lg whitespace-pre-wrap break-words border border-[#A93E39] h-[500px] overflow-auto"
        >
          {output}
        </pre>
      </form>
    </main>
  );
}
