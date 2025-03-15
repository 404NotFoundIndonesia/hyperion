
import { 
  openFiles, 
  openFolder, 
  exportFiles, 
  removeAllFiles, 
  unloadSelectedFile,
  obfuscateAll,
  obfuscateSelectedFile,
  toggleConfig,
} from "./appActions.js";

export const menuConfig = [
    {
        name: "File",
        children: [
          { name: "Open", shortcut: "Ctrl+O", action: async () => await openFiles() }, // Ensure async execution
          { name: "Open Folder", shortcut: "Ctrl+Shift+O", action: async () => await openFolder()},
          { name: "Save", shortcut: "Ctrl+S", action: () => console.log("Save") },
          { name: "Save As", shortcut: "Ctrl+Shift+S", action: async () => await exportFiles() },
          { name: "Close", shortcut: "Ctrl+W", action: () => unloadSelectedFile() },
          { name: "Close Folder", shortcut: "Ctrl+Shift+W", action: () => removeAllFiles() }
        ]
      },
    {
      name: "Run",
      children: [
        { name: "Obfuscate", action: async () => await obfuscateSelectedFile() },
        { name: "Obfuscate All", action: async () => await obfuscateAll() },
        { name: "Configuration", action: () => toggleConfig() },
      ]
    },
    {
      name: "Help",
      children: [
        { name: "How to Contribute", action: () => console.log("How to Contribute") },
        { name: "Report Issue", action: () => console.log("Report Issue") }
      ]
    }
  ];
  