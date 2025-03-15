import {
  selectedFile,
  selectedFiles,
  fileTree,
  filesContent,
  obfuscatedContent,
  showContent,
} from "./stores/selectedFileStore.js";

import {
  OpenFiles,
  OpenFolder,
  ReadFilesContent,
  ObfuscateJS,
  ExportEncryptedFiles,
  OpenImportedFileLocation,
} from "../wailsjs/go/main/App";

import { get } from "svelte/store";
import { obfuscationConfig } from "./stores/configStore";

export async function openFiles() {
  try {
    const newFiles = await OpenFiles();
    if (newFiles.length > 0) {
      selectedFiles.update((files) => [
        ...new Set([...(files || []), ...newFiles]),
      ]);
      fileTree.set(buildFileTree(get(selectedFiles)));
    }
  } catch (error) {
    console.error("Error opening file dialog:", error);
  }
}

export async function openFolder() {
  try {
    const folderFiles = await OpenFolder();
    if (folderFiles && folderFiles.length > 0) {
      selectedFiles.update((files) => {
        const updatedFiles = [...new Set([...(files || []), ...folderFiles])];
        fileTree.set(buildFileTree(updatedFiles));
        return updatedFiles;
      });
    }
  } catch (error) {
    console.error("Error opening folder dialog:", error);
  }
}

export function buildFileTree(files) {
  const tree = {};

  files.forEach((file) => {
    const normalizedFile = file.replace(/\\/g, "/");
    const parts = normalizedFile.split("/");

    let current = tree;
    parts.forEach((part, index) => {
      if (!current[part]) {
        current[part] =
          index === parts.length - 1 ? { path: normalizedFile } : {};
      }
      current = current[part];
    });
  });

  let keys = Object.keys(tree);
  let trimmedTree = tree;

  while (keys.length === 1) {
    trimmedTree = trimmedTree[keys[0]];
    keys = Object.keys(trimmedTree);
  }

  if (typeof trimmedTree !== "object" || Array.isArray(trimmedTree)) {
    const path = trimmedTree;
    const fileName = trimmedTree.split("/").at(-1);
    console.log(fileName);
    trimmedTree = { [fileName]: { path: path } };
  }

  console.log(trimmedTree);
  return trimmedTree;
}

export async function exportFiles() {
  if (
    get(obfuscatedContent) &&
    Object.keys(get(obfuscatedContent)).length === 0
  ) {
    console.error("No encrypted files to export!");
    return;
  }

  try {
    // @ts-ignore
    await ExportEncryptedFiles(get(obfuscatedContent));
    console.log("Files exported successfully!");
  } catch (error) {
    console.error("Error exporting files:", error);
  }
}

export function removeAllFiles() {
  selectedFiles.set([]);
  filesContent.set({});
  obfuscatedContent.set({});
  showContent.set({});
  selectedFile.set(null);
}

export function unloadFile(filePath) {
  filesContent.update((content) => {
    if (content[filePath]) {
      const newContent = { ...content };
      delete newContent[filePath];
      return newContent;
    }
    return content;
  });

  if (get(selectedFile) === filePath) {
    selectedFile.set(null);
  }
}

export function unloadSelectedFile() {
  filesContent.update((content) => {
    if (content[get(selectedFile)]) {
      const newContent = { ...content };
      delete newContent[get(selectedFile)];
      return newContent;
    }
    return content;
  });

  if (get(selectedFile) === get(selectedFile)) {
    selectedFile.set(null);
  }
}

export async function obfuscateAll() {
  if (get(selectedFiles).length === 0) {
    console.warn("No files selected!");
    return;
  }

  try {
    const newFilesContent = await ReadFilesContent(get(selectedFiles));
    filesContent.update((fc) => ({ ...fc, ...newFilesContent }));
    const updatedObfuscatedContent = { ...get(obfuscatedContent) };
    for (let filePath of Object.keys(newFilesContent)) {
      updatedObfuscatedContent[filePath] = await ObfuscateJS(
        newFilesContent[filePath],
        { ...get(obfuscationConfig) },
      );
    }
    obfuscatedContent.set(updatedObfuscatedContent);
    const filePaths = Object.keys(get(filesContent));
    if (filePaths.length > 0 && !get(selectedFile)) {
      selectedFile.set(filePaths[0]);
    }
  } catch (error) {
    console.error("Error processing files:", error);
  }
}

export async function obfuscateFile(filePath) {
  if (!filePath || !get(filesContent)[filePath]) {
    console.warn("Invalid file or file not found!");
    return;
  }

  try {
    const obfuscated = await ObfuscateJS(get(filesContent)[filePath], {
      ...get(obfuscationConfig),
    });

    obfuscatedContent.update((oc) => ({ ...oc, [filePath]: obfuscated }));
  } catch (error) {
    console.error("Error obfuscating file:", error);
  }
}

export async function obfuscateSelectedFile() {
  const filePath = get(selectedFile);
  if (!filePath || !get(filesContent)[filePath]) {
    console.warn("Invalid file or file not found!");
    return;
  }
  try {
    const obfuscated = await ObfuscateJS(get(filesContent)[filePath], {
      ...get(obfuscationConfig),
    });
    obfuscatedContent.update((oc) => ({ ...oc, [filePath]: obfuscated }));
  } catch (error) {
    console.error("Error obfuscating file:", error);
  }
}

export function toggleConfig() {
  document.getElementById("config-container").classList.toggle("hidden");
}