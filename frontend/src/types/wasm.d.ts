export interface Go {
  importObject: any;
  run(instance: WebAssembly.Instance): Promise<void>;
}

declare global {
  interface Window {
    Go: new () => Go;
  }
}

// Ensure the file is treated as a module
export {};
