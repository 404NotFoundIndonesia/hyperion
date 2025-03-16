<script>
  // @ts-nocheck

  import {
    OpenFiles,
    OpenFolder,
    ReadFilesContent,
    ObfuscateJS,
    ExportEncryptedFiles,
    OpenImportedFileLocation,
  } from "../wailsjs/go/main/App";

  import logo from "/src/assets/images/hyperion-blue.png";
  import { writable, get } from "svelte/store";
  import { obfuscationConfig } from "./stores/configStore";
  import { onMount, onDestroy } from "svelte";
  import { menuConfig } from "./menuConfig.js";
  import {
    selectedFile,
    selectedFiles,
    fileTree,
    filesContent,
    obfuscatedContent,
    showContent,
  } from "./stores/selectedFileStore.js";
  import {
    openFiles,
    openFolder,
    buildFileTree,
    exportFiles,
    removeAllFiles,
    unloadFile,
    obfuscateAll,
    obfuscateFile,
    toggleConfig,
  } from "./appActions.js";
  import { sidebarWidth, collapsedFolders } from "./stores/uiStore";

  let previewOriginal = writable(true);
  let isMac = navigator.userAgent.includes("Mac");
  let isFrameless = false;
  // isMac = true;
  let isFullscreen = false;
  let isOpen = true;
  let buttons = [];
  let menuContainer;
  let isMaximize = false;
  let activeMenu = null;
  let isResizing = false;

  function startResizing(event) {
    isResizing = true;
    const startX = event.clientX;
    const startWidth = $sidebarWidth;

    function resize(event) {
      if (isResizing) {
        const newWidth = startWidth + (event.clientX - startX);
        sidebarWidth.set(Math.max(150, newWidth)); // Set a min width
      }
    }

    function stopResizing() {
      isResizing = false;
      window.removeEventListener("mousemove", resize);
      window.removeEventListener("mouseup", stopResizing);
    }

    window.addEventListener("mousemove", resize);
    window.addEventListener("mouseup", stopResizing);
  }

  function toggleFolder(name) {
    collapsedFolders.update((state) => ({
      ...state,
      [name]: !state[name],
    }));
  }

  function toggleMenu(name) {
    activeMenu = activeMenu === name ? null : name;
  }

  function executeAction(action) {
    if (typeof action === "function") action();
    setTimeout(() => {
      activeMenu = null;
    }, 0);
  }

  function handleClickOutside(event) {
    const menuElement = document.getElementById("windowsMenu");
    if (menuElement && !menuElement.contains(event.target)) {
      activeMenu = null;
    }
  }

  sidebarWidth.subscribe((value) => {
    localStorage.setItem("sidebarWidth", value);
  });

  onMount(() => {
    document.addEventListener("click", handleClickOutside);
  });

  onDestroy(() => {
    document.removeEventListener("click", handleClickOutside);
  });

  fileTree.set(buildFileTree($selectedFiles));

  async function selectFile(filePath) {
    if ($selectedFile === filePath) return;

    selectedFile.set(filePath);
    fileTree.set(buildFileTree($selectedFiles));
    if (!get(filesContent)[filePath]) {
      const content = await ReadFilesContent([filePath]);

      filesContent.update((fc) => ({ ...fc, ...content }));
    }

    setTimeout(() => {
      const filteredIndex = $selectedFiles
        .filter((file) => filesContent[file])
        .indexOf(filePath);

      if (buttons[filteredIndex]) {
        buttons[filteredIndex].scrollIntoView({
          behavior: "smooth",
          block: "nearest",
        });
      }
    }, 50);
  }

  function handleClick(event) {
    console.log("handleClick triggered", event.target);

    const target = event.target.closest("[data-value], [data-folder]");

    if (target) {
      const folderName = target.getAttribute("data-folder");
      
      if (folderName) {
        console.log("Folder clicked:", folderName);
        collapsedFolders.update((state) => {
          const newState = { ...state, [folderName]: !state[folderName] };
          console.log("Updated collapsedFolders:", newState);
          return newState;
        });

        console.log("Collapsed state after update:", get(collapsedFolders));
      } else {
        let value = target.getAttribute("data-value");
        value = value.replace(/\//g, "\\");
        console.log("File selected:", value);
        selectFile(value);
      }
      fileTree.set(buildFileTree($selectedFiles));
    }
  }

  function renderTree(tree) {
    console.log(
      "Rendering tree with collapsed folders:",
      get(collapsedFolders),
    );

    return Object.entries(tree)
      .map(([name, content]) => {
        const isLeaf = !(content && typeof content === "object");
        if (!isLeaf) {
          const isCollapsed = get(collapsedFolders)[name] ?? false;
          const firstChild = Object.entries(content)[0];
          const isChildLeaf =
            firstChild && !(firstChild[1] && typeof firstChild[1] === "object");
          let isActive = false;
          if (isChildLeaf) {
            isActive = firstChild[1].replace(/\//g, "\\") == get(selectedFile);
          }
          return isChildLeaf
            ? `<div data-value="${firstChild[1]}" class="${isActive ? "opacity-100" : "opacity-50"} cursor-pointer pl-2 hover:opacity-100 flex gap-1 items-center">
              <svg  xmlns="http://www.w3.org/2000/svg"  width="14"  height="14"  viewBox="0 0 24 24"  fill="none"  stroke="currentColor"  stroke-width="2"  stroke-linecap="round"  stroke-linejoin="round"  class="shrink-0 icon icon-tabler icons-tabler-outline icon-tabler-file"><path stroke="none" d="M0 0h24v24H0z" fill="none"/><path d="M14 3v4a1 1 0 0 0 1 1h4" /><path d="M17 21h-10a2 2 0 0 1 -2 -2v-14a2 2 0 0 1 2 -2h7l5 5v11a2 2 0 0 1 -2 2z" /></svg>
              ${name}
              </div>`
            : `
          <ul class="list-none pl-2">
            <li class="pl-2 dark:text-white">
              <div data-folder="${name}" class="flex gap-1 items-center cursor-pointer">
                <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="shrink-0 icon icon-tabler icons-tabler-outline icon-tabler-folder">
                  <path stroke="none" d="M0 0h24v24H0z" fill="none"/>
                  <path d="M5 4h4l3 3h7a2 2 0 0 1 2 2v8a2 2 0 0 1 -2 2h-14a2 2 0 0 1 -2 -2v-11a2 2 0 0 1 2 -2" />
                </svg>
                ${name}
              </div>
              ${isCollapsed ? "" : renderTree(content)}
            </li>
          </ul>`;
        }
      })
      .join("");
  }

  function toggleSidebar() {
    isOpen = !isOpen;
  }

  async function checkFullscreen() {
    isFullscreen = await window.runtime.WindowIsFullscreen();
    isMaximize = await window.runtime.WindowIsMaximised();
  }

  window.addEventListener("resize", checkFullscreen);

  setTimeout(() => {
    checkFullscreen();
  }, 500);

  function togglePreview() {
    previewOriginal.update((prev) => !prev);
  }

  function removeFile(index) {
    selectedFiles.update((files) => {
      if (!Array.isArray(files) || index < 0 || index >= files.length)
        return files;

      const removedFile = files[index];
      const updatedFiles = files.filter((_, i) => i !== index);

      delete filesContent[removedFile];
      delete obfuscatedContent[removedFile];
      delete showContent[removedFile];

      if (get(selectedFile) === removedFile) {
        const filePaths = Object.keys(filesContent);
        selectedFile.set(filePaths.length > 0 ? filePaths[0] : null);
      }

      return updatedFiles;
    });
  }

  function getFileName(path) {
    return path.split(/[\\/]/).pop();
  }

  function toggleShowCode(filePath) {
    showContent.update((sc) => ({
      ...sc,
      [filePath]: !sc[filePath],
    }));
  }

  async function copyOutput(content) {
    try {
      await navigator.clipboard.writeText(content);
    } catch (err) {
      console.error("Failed to copy:", err);
    }
  }

  async function openImportedFolder(filePath) {
    if (filePath) {
      await OpenImportedFileLocation(filePath);
    }
  }
</script>

<div id="config-container" class="fixed p-10 h-full w-full z-1 flex hidden">
  <!-- svelte-ignore a11y_no_static_element_interactions -->
  <!-- svelte-ignore a11y_click_events_have_key_events -->
  <div on:click={toggleConfig} class="fixed h-full w-full top-0 left-0"></div>
  <div
    class="relative bg-neutral-100 dark:bg-neutral-800 w-full max-w-screen-lg mx-auto h-full border border-black/15 dark:border-white/15 rounded-lg shadow-xl"
  >
    <div
      class="flex flex-col gap-1 p-2 overflow-auto mt-[40px] h-[calc(100%-40px)]"
    >
      <div
        style="--wails-draggable:drag;"
        class="absolute top-0 left-0 w-full flex gap-1 items-center dark:text-white border-b border-black/15 dark:border-white/15"
      >
        <!-- svelte-ignore a11y_no_static_element_interactions -->
        <!-- svelte-ignore a11y_click_events_have_key_events -->
        <div
          on:click={toggleConfig}
          style="--wails-draggable: false"
          class="dark:text-white p-2 ms-1 cursor-pointer {isMac
            ? ''
            : 'hidden'}"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="12"
            height="12"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2.5"
            stroke-linecap="round"
            stroke-linejoin="round"
            class="icon icon-tabler icons-tabler-outline icon-tabler-x text-red-600 hover:text-red-900 bg-red-600 rounded-full"
            ><path stroke="none" d="M0 0h24v24H0z" fill="none" /><path
              d="M18 6l-12 12"
            /><path d="M6 6l12 12" /></svg
          >
        </div>
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="16"
          height="16"
          viewBox="0 0 24 24"
          fill="currentColor"
          class="icon icon-tabler icons-tabler-filled icon-tabler-settings text-blue-400 {isMac
            ? 'ms-auto hidden'
            : 'ms-2'}"
          ><path stroke="none" d="M0 0h24v24H0z" fill="none" /><path
            d="M14.647 4.081a.724 .724 0 0 0 1.08 .448c2.439 -1.485 5.23 1.305 3.745 3.744a.724 .724 0 0 0 .447 1.08c2.775 .673 2.775 4.62 0 5.294a.724 .724 0 0 0 -.448 1.08c1.485 2.439 -1.305 5.23 -3.744 3.745a.724 .724 0 0 0 -1.08 .447c-.673 2.775 -4.62 2.775 -5.294 0a.724 .724 0 0 0 -1.08 -.448c-2.439 1.485 -5.23 -1.305 -3.745 -3.744a.724 .724 0 0 0 -.447 -1.08c-2.775 -.673 -2.775 -4.62 0 -5.294a.724 .724 0 0 0 .448 -1.08c-1.485 -2.439 1.305 -5.23 3.744 -3.745a.722 .722 0 0 0 1.08 -.447c.673 -2.775 4.62 -2.775 5.294 0zm-2.647 4.919a3 3 0 1 0 0 6a3 3 0 0 0 0 -6z"
          /></svg
        >
        <div class="text-xs {isMac ? 'ms-auto me-8' : ''}">
          Obfuscator Config
        </div>
        <!-- svelte-ignore a11y_no_static_element_interactions -->
        <!-- svelte-ignore a11y_click_events_have_key_events -->
        <div
          on:click={toggleConfig}
          style="--wails-draggable: false"
          class="ms-auto dark:text-white p-2 cursor-pointer hover:bg-red-700 hover:text-white hover:text-white rounded-tr-lg"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="18"
            height="18"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="1"
            stroke-linecap="round"
            stroke-linejoin="round"
            class="icon icon-tabler icons-tabler-outline icon-tabler-x {isMac
              ? 'hidden'
              : ''}"
            ><path stroke="none" d="M0 0h24v24H0z" fill="none" /><path
              d="M18 6l-12 12"
            /><path d="M6 6l12 12" /></svg
          >
        </div>
      </div>
      <label class="input-setting">
        <input
          type="checkbox"
          bind:checked={$obfuscationConfig.removeWhiteSpace}
        />
        Remove Whitespace
      </label>

      <label class="input-setting">
        <input
          type="checkbox"
          bind:checked={$obfuscationConfig.removeComments}
        />
        Remove Comments
      </label>

      <label class="input-setting">
        <input
          type="checkbox"
          bind:checked={$obfuscationConfig.changeVariable}
        />
        Change Variable Names
      </label>
    </div>
  </div>
</div>
<div
  class="flex flex-col h-full transition-colors duration-300 overflow-hidden {isFullscreen
    ? 'bg-white dark:bg-neutral-900'
    : isMac
      ? 'bg-white/50 dark:bg-neutral-900/50'
      : 'border-t border-black/15 dark:border-white/15'}"
>
  <div
    id="topbar"
    style="--wails-draggable:drag;"
    class="w-full border-b border-black/15 dark:border-white/15 shrink-0"
  >
    <div class="h-[32px] flex gap-1 shrink-0 items-center">
      <img
        class="ms-2 h-[16px] {isMac || !isFrameless ? 'hidden' : ''}"
        src={logo}
        alt="app logo"
      />
      <div id="windowsMenu" class="flex h-full {isMac ? 'hidden' : ''}">
        {#each menuConfig as menu}
          <div class="relative">
            <!-- svelte-ignore a11y_no_static_element_interactions -->
            <!-- svelte-ignore a11y_click_events_have_key_events -->
            <div class="menuButton" on:click={() => toggleMenu(menu.name)}>
              {menu.name}
            </div>

            {#if activeMenu === menu.name}
              <div class="menuDropdown z-10">
                {#each menu.children as item}
                  <!-- svelte-ignore a11y_click_events_have_key_events -->
                  <!-- svelte-ignore a11y_no_static_element_interactions -->
                  <div
                    on:click={() => executeAction(item.action)}
                    class="menuItem flex gap-1 items-center justify-between"
                  >
                    <div class="text-sm">
                      {item.name}
                    </div>
                    <div class="text-xs">
                      {item.shortcut ? `(${item.shortcut})` : ""}
                    </div>
                  </div>
                {/each}
              </div>
            {/if}
          </div>
        {/each}
      </div>
      <div class="flex gap-1 p-1 ms-auto">
        <button
          class="btn-sm"
          on:click={obfuscateAll}
          disabled={$selectedFiles.length === 0}
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="14"
            height="14"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
            class="icon icon-tabler icons-tabler-outline icon-tabler-shield-code"
            ><path stroke="none" d="M0 0h24v24H0z" fill="none" /><path
              d="M12 21a12 12 0 0 1 -8.5 -15a12 12 0 0 0 8.5 -3a12 12 0 0 0 8.5 3a12 12 0 0 1 -.078 7.024"
            /><path d="M20 21l2 -2l-2 -2" /><path d="M17 17l-2 2l2 2" /></svg
          >
          Obfuscate All
        </button>
        <button class="btn-sm" on:click={exportFiles} disabled="true">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="14"
            height="14"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
            class="icon icon-tabler icons-tabler-outline icon-tabler-device-floppy"
            ><path stroke="none" d="M0 0h24v24H0z" fill="none" /><path
              d="M6 4h10l4 4v10a2 2 0 0 1 -2 2h-12a2 2 0 0 1 -2 -2v-12a2 2 0 0 1 2 -2"
            /><path d="M12 14m-2 0a2 2 0 1 0 4 0a2 2 0 1 0 -4 0" /><path
              d="M14 4l0 4l-6 0l0 -4"
            /></svg
          >
          Save
        </button>
        <button
          class="btn-sm"
          on:click={exportFiles}
          disabled={Object.keys($obfuscatedContent).length === 0}
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="14"
            height="14"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
            class="icon icon-tabler icons-tabler-outline icon-tabler-device-floppy"
            ><path stroke="none" d="M0 0h24v24H0z" fill="none" /><path
              d="M6 4h10l4 4v10a2 2 0 0 1 -2 2h-12a2 2 0 0 1 -2 -2v-12a2 2 0 0 1 2 -2"
            /><path d="M12 14m-2 0a2 2 0 1 0 4 0a2 2 0 1 0 -4 0" /><path
              d="M14 4l0 4l-6 0l0 -4"
            /></svg
          >
          Save As
        </button>
      </div>
      <div
        id="windowsControl"
        class="flex h-full {isMac || !isFrameless ? 'hidden' : ''}"
      >
        <!-- svelte-ignore a11y_click_events_have_key_events -->
        <!-- svelte-ignore a11y_no_static_element_interactions -->
        <div
          class="windowsButton"
          on:click={() => {
            window.runtime.WindowMinimise();
          }}
        >
          <svg
            width="10"
            height="1"
            viewBox="0 0 10 1"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <rect width="10" height="1" fill="white" />
          </svg>
        </div>
        <!-- svelte-ignore a11y_click_events_have_key_events -->
        <!-- svelte-ignore a11y_no_static_element_interactions -->
        <div
          class="windowsButton"
          on:click={() => {
            window.runtime.WindowToggleMaximise();
          }}
        >
          <svg
            class={isMaximize ? "hidden" : ""}
            width="11"
            height="10"
            viewBox="0 0 11 10"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              d="M8 0.5H3C1.89543 0.5 1 1.39543 1 2.5V7.5C1 8.60457 1.89543 9.5 3 9.5H8C9.10457 9.5 10 8.60457 10 7.5V2.5C10 1.39543 9.10457 0.5 8 0.5Z"
              stroke="white"
            />
          </svg>
          <svg
            class={isMaximize ? "" : "hidden"}
            width="11"
            height="10"
            viewBox="0 0 11 10"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              d="M2.5 0.5H8C9.10457 0.5 10 1.39543 10 2.5V8M3 9.5H6C7.10457 9.5 8 8.60457 8 7.5V4.5C8 3.39543 7.10457 2.5 6 2.5H3C1.89543 2.5 1 3.39543 1 4.5V7.5C1 8.60457 1.89543 9.5 3 9.5Z"
              stroke="white"
            />
          </svg>
        </div>
        <!-- svelte-ignore a11y_click_events_have_key_events -->
        <!-- svelte-ignore a11y_no_static_element_interactions -->
        <div
          class="windowsButton closeButton"
          on:click={() => {
            window.runtime.Quit();
          }}
        >
          <svg
            width="12"
            height="12"
            viewBox="0 0 12 12"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path d="M1 11L6 6M11 1L6 6M6 6L11 11M6 6L1 1" stroke="white" />
          </svg>
        </div>
      </div>
    </div>
  </div>
  <div
    style="height: {isMac || !isFrameless ? 'calc(100% - 32px);' : ''}"
    class="flex h-full"
  >
    <div
      style="--wails-draggable:drag;"
      class="flex flex-col border-r border-black/15 dark:border-white/15 bg-black/10"
    >
      <!-- svelte-ignore a11y_consider_explicit_label -->
      <button
        disabled={$selectedFiles.length === 0}
        class="text-black dark:text-white disabled:opacity-50 p-2 enabled:cursor-pointer flex gap-1 items-center {$selectedFiles.length >
          0 && isOpen
          ? 'bg-white dark:bg-white/10'
          : ''}"
        on:click={toggleSidebar}
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="30"
          height="30"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="1.5"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="icon icon-tabler icons-tabler-outline icon-tabler-files"
          ><path stroke="none" d="M0 0h24v24H0z" fill="none" /><path
            d="M15 3v4a1 1 0 0 0 1 1h4"
          /><path
            d="M18 17h-7a2 2 0 0 1 -2 -2v-10a2 2 0 0 1 2 -2h4l5 5v7a2 2 0 0 1 -2 2z"
          /><path
            d="M16 17v2a2 2 0 0 1 -2 2h-7a2 2 0 0 1 -2 -2v-10a2 2 0 0 1 2 -2h2"
          /></svg
        >
      </button>
      <!-- svelte-ignore a11y_consider_explicit_label -->
      <button
        class="text-black dark:text-white disabled:opacity-50 p-2 enabled:cursor-pointer flex gap-1 items-center"
        on:click={openFiles}
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="30"
          height="30"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="1.5"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="icon icon-tabler icons-tabler-outline icon-tabler-file-code"
          ><path stroke="none" d="M0 0h24v24H0z" fill="none" /><path
            d="M14 3v4a1 1 0 0 0 1 1h4"
          /><path
            d="M17 21h-10a2 2 0 0 1 -2 -2v-14a2 2 0 0 1 2 -2h7l5 5v11a2 2 0 0 1 -2 2z"
          /><path d="M10 13l-1 2l1 2" /><path d="M14 13l1 2l-1 2" /></svg
        >
      </button>
      <!-- svelte-ignore a11y_consider_explicit_label -->
      <button
        class="text-black dark:text-white disabled:opacity-50 p-2 enabled:cursor-pointer flex gap-1 items-center"
        on:click={openFolder}
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="30"
          height="30"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="1.5"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="icon icon-tabler icons-tabler-outline icon-tabler-folder-code"
          ><path stroke="none" d="M0 0h24v24H0z" fill="none" /><path
            d="M11 19h-6a2 2 0 0 1 -2 -2v-11a2 2 0 0 1 2 -2h4l3 3h7a2 2 0 0 1 2 2v4"
          /><path d="M20 21l2 -2l-2 -2" /><path d="M17 17l-2 2l2 2" /></svg
        >
      </button>
      <!-- svelte-ignore a11y_consider_explicit_label -->
      <button
        class="text-black dark:text-white disabled:opacity-50 p-2 enabled:cursor-pointer flex gap-1 items-center mt-auto"
        on:click={toggleConfig}
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="30"
          height="30"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="1.5"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="icon icon-tabler icons-tabler-outline icon-tabler-settings"
          ><path stroke="none" d="M0 0h24v24H0z" fill="none" /><path
            d="M10.325 4.317c.426 -1.756 2.924 -1.756 3.35 0a1.724 1.724 0 0 0 2.573 1.066c1.543 -.94 3.31 .826 2.37 2.37a1.724 1.724 0 0 0 1.065 2.572c1.756 .426 1.756 2.924 0 3.35a1.724 1.724 0 0 0 -1.066 2.573c.94 1.543 -.826 3.31 -2.37 2.37a1.724 1.724 0 0 0 -2.572 1.065c-.426 1.756 -2.924 1.756 -3.35 0a1.724 1.724 0 0 0 -2.573 -1.066c-1.543 .94 -3.31 -.826 -2.37 -2.37a1.724 1.724 0 0 0 -1.065 -2.572c-1.756 -.426 -1.756 -2.924 0 -3.35a1.724 1.724 0 0 0 1.066 -2.573c-.94 -1.543 .826 -3.31 2.37 -2.37c1 .608 2.296 .07 2.572 -1.065z"
          /><path d="M9 12a3 3 0 1 0 6 0a3 3 0 0 0 -6 0" /></svg
        >
      </button>
    </div>
    {#if $selectedFiles.length === 0}
      <div
        class="flex flex-grow w-full p-20 overflow-hidden dark:text-white items-center max-w-screen-md mx-auto"
      >
        <div class="flex flex-col w-full gap-2">
          <div class="text-2xl font-medium">Hyperion</div>
          <div class="opacity-50 font-light">Fast and Secure</div>
          <div class="text-lg font-medium">Start</div>
          <button class="btn me-auto w-32" on:click={openFiles}>
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="18"
              height="18"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="1.5"
              stroke-linecap="round"
              stroke-linejoin="round"
              class="icon icon-tabler icons-tabler-outline icon-tabler-file-code"
              ><path stroke="none" d="M0 0h24v24H0z" fill="none" /><path
                d="M14 3v4a1 1 0 0 0 1 1h4"
              /><path
                d="M17 21h-10a2 2 0 0 1 -2 -2v-14a2 2 0 0 1 2 -2h7l5 5v11a2 2 0 0 1 -2 2z"
              /><path d="M10 13l-1 2l1 2" /><path d="M14 13l1 2l-1 2" /></svg
            >
            Open File
          </button>
          <button class="btn me-auto w-32" on:click={openFolder}>
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="18"
              height="18"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="1.5"
              stroke-linecap="round"
              stroke-linejoin="round"
              class="icon icon-tabler icons-tabler-outline icon-tabler-folder-code"
              ><path stroke="none" d="M0 0h24v24H0z" fill="none" /><path
                d="M11 19h-6a2 2 0 0 1 -2 -2v-11a2 2 0 0 1 2 -2h4l3 3h7a2 2 0 0 1 2 2v4"
              /><path d="M20 21l2 -2l-2 -2" /><path d="M17 17l-2 2l2 2" /></svg
            >
            Open Folder
          </button>
          <div class="text-lg font-medium">Contribute</div>
          <button
            class="btn me-auto"
            on:click={() => {
              window.runtime.BrowserOpenURL(
                "https://github.com/404NotFoundIndonesia/hyperion",
              );
            }}
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="18"
              height="18"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="1.5"
              stroke-linecap="round"
              stroke-linejoin="round"
              class="icon icon-tabler icons-tabler-outline icon-tabler-brand-github"
              ><path stroke="none" d="M0 0h24v24H0z" fill="none" /><path
                d="M9 19c-4.3 1.4 -4.3 -2.5 -6 -3m12 5v-3.5c0 -1 .1 -1.4 -.5 -2c2.8 -.3 5.5 -1.4 5.5 -6a4.6 4.6 0 0 0 -1.3 -3.2a4.2 4.2 0 0 0 -.1 -3.2s-1.1 -.3 -3.5 1.3a12.3 12.3 0 0 0 -6.2 0c-2.4 -1.6 -3.5 -1.3 -3.5 -1.3a4.2 4.2 0 0 0 -.1 3.2a4.6 4.6 0 0 0 -1.3 3.2c0 4.6 2.7 5.7 5.5 6c-.6 .6 -.6 1.2 -.5 2v3.5"
              /></svg
            >
            Open Repository
          </button>
        </div>
        <img
          class="min-w-1 mx-auto my-auto object-scale-down max-h-full max-h-full grayscale opacity-25 pointer-events-none"
          src={logo}
          alt="app logo"
        />
      </div>
    {/if}
    <div
      class="flex flex overflow-hidden grow {$selectedFiles.length == 0
        ? 'hidden'
        : ''}"
    >
      {#if $selectedFiles.length > 0 && isOpen}
        <div
          id="sidebar"
          style="width:{$sidebarWidth}px"
          class="flex relative dark:text-white bg-black/5 dark:bg-white/5 shrink-0"
        >
          <div class="w-full pe-2">
            <div class="h-[40px] flex items-center p-1 px-2 pe-0.5 text-sm">
              Imported Files
              <!-- svelte-ignore a11y_consider_explicit_label -->
              <div class="ms-auto">
                <button
                  class="text-black dark:text-white disabled:opacity-50 p-2 enabled:cursor-pointer flex gap-1 items-center me-[-6px]"
                  on:click={removeAllFiles}
                  disabled={$selectedFiles.length === 0}
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="20"
                    height="20"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="1.5"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    class="icon icon-tabler icons-tabler-outline icon-tabler-copy-minus"
                    ><path stroke="none" d="M0 0h24v24H0z" fill="none" /><path
                      stroke="none"
                      d="M0 0h24v24H0z"
                    /><path
                      d="M7 9.667a2.667 2.667 0 0 1 2.667 -2.667h8.666a2.667 2.667 0 0 1 2.667 2.667v8.666a2.667 2.667 0 0 1 -2.667 2.667h-8.666a2.667 2.667 0 0 1 -2.667 -2.667z"
                    /><path
                      d="M4.012 16.737a2 2 0 0 1 -1.012 -1.737v-10c0 -1.1 .9 -2 2 -2h10c.75 0 1.158 .385 1.5 1"
                    /><path d="M11 14h6" /></svg
                  >
                </button>
              </div>
            </div>
            <div
              class="sidebar-scroll text-sm overflow-x-auto overflow-y-auto h-[calc(100%-40px)]"
              style=""
            >
              <ul
                style="width: 100%"
                class="flex flex-col min-h-full shrink-0 dark:text-white"
              >
                <!-- svelte-ignore a11y_click_events_have_key_events -->
                <!-- svelte-ignore a11y_no_static_element_interactions -->
                <div on:click={handleClick}>
                  {@html renderTree($fileTree)}
                </div>
              </ul>
            </div>
          </div>
          <!-- svelte-ignore a11y_no_static_element_interactions -->
          <div
            class="resizer w-2 border-r border-black/10 dark:border-white/10 h-full absolute right-0 cursor-ew-resize"
            on:mousedown={startResizing}
          ></div>
        </div>
      {/if}
      <div class="flex flex-col overflow-x-hidden overflow-y-auto grow">
        <div
          class="shrink-0 h-[40px] flex gap-1 dark:text-white w-full border-b border-black/15 dark:border-white/15"
        >
          <div
            class="grow text-sm overflow-x-auto"
            style="scrollbar-width: none;"
            bind:this={menuContainer}
          >
            <div class="flex flex-row h-full">
              {#each $selectedFiles.filter((file) => $filesContent[file]) as filePath, i}
                <button
                  class="btn-block"
                  class:active={$selectedFile === filePath}
                  bind:this={buttons[i]}
                >
                  <div class="flex w-full">
                    <!-- svelte-ignore a11y_click_events_have_key_events -->
                    <!-- svelte-ignore a11y_no_static_element_interactions -->
                    <div
                      class="flex flex-row w-full items-center grow ps-3 overflow-hidden"
                      on:click={() => selectFile(filePath)}
                    >
                      {#if $obfuscatedContent[filePath]}
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          width="16"
                          height="16"
                          viewBox="0 0 24 24"
                          fill="none"
                          stroke="currentColor"
                          stroke-width="2"
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          class="shrink-0 icon icon-tabler icons-tabler-outline icon-tabler-lock"
                        >
                          <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                          <path
                            d="M5 13a2 2 0 0 1 2 -2h10a2 2 0 0 1 2 2v6a2 2 0 0 1 -2 2h-10a2 2 0 0 1 -2 -2v-6z"
                          />
                          <path d="M11 16a1 1 0 1 0 2 0a1 1 0 0 0 -2 0" />
                          <path d="M8 11v-4a4 4 0 1 1 8 0v4" />
                        </svg>
                      {:else}
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          width="16"
                          height="16"
                          viewBox="0 0 24 24"
                          fill="none"
                          stroke="currentColor"
                          stroke-width="2"
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          class="shrink-0 icon icon-tabler icons-tabler-outline icon-tabler-lock-open"
                        >
                          <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                          <path
                            d="M5 11m0 2a2 2 0 0 1 2 -2h10a2 2 0 0 1 2 2v6a2 2 0 0 1 -2 2h-10a2 2 0 0 1 -2 -2z"
                          />
                          <path d="M12 16m-1 0a1 1 0 1 0 2 0a1 1 0 1 0 -2 0" />
                          <path d="M8 11v-5a4 4 0 0 1 8 0" />
                        </svg>
                      {/if}
                      <div class="grow overflow-hidden">
                        <div class="whitespace-nowrap grow">
                          {getFileName(filePath)}
                        </div>
                      </div>
                    </div>
                    <!-- svelte-ignore a11y_click_events_have_key_events -->
                    <!-- svelte-ignore a11y_no_static_element_interactions -->
                    <div
                      class="flex items-center px-1 pe-2 shrink-0 opacity-50 hover:opacity-100 hover:text-red-700"
                      on:click={() => unloadFile(filePath)}
                    >
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width="16"
                        height="16"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        class="icon icon-tabler icons-tabler-outline icon-tabler-x"
                        ><path
                          stroke="none"
                          d="M0 0h24v24H0z"
                          fill="none"
                        /><path d="M18 6l-12 12" /><path d="M6 6l12 12" /></svg
                      >
                    </div>
                  </div>
                </button>
              {/each}
            </div>
          </div>
        </div>
        {#if Object.keys($filesContent).length > 0}
          <div style="height: 200px;" class="flex flex-col gap-2 p-2 grow">
            {#if $selectedFile}
              <div
                class="max-h-full flex flex-col bg-white dark:bg-neutral-700/50 border border-black/15 dark:border-white/15 dark:text-white rounded-md"
              >
                <div
                  class="p-4 flex flex-col md:flex-row gap-2 w-full md:items-center border-b border-black/15 dark:border-white/15"
                >
                  <div class="shrink">
                    {getFileName($selectedFile)}
                    <div
                      class="text-xs text-neutral-500 whitespace-pre-wrap break-all"
                    >
                      {$selectedFile}
                    </div>
                  </div>
                  {#if $obfuscatedContent[$selectedFile]}
                    <div class="flex ms-auto gap-1 shrink-0 items-start">
                      <button
                        class="btn"
                        on:click={() => toggleShowCode($selectedFile)}
                      >
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          width="18"
                          height="18"
                          viewBox="0 0 24 24"
                          fill="none"
                          stroke="currentColor"
                          stroke-width="2"
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          class="icon icon-tabler icons-tabler-outline icon-tabler-eye"
                        >
                          <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                          <path d="M10 12a2 2 0 1 0 4 0a2 2 0 0 0 -4 0" />
                          <path
                            d="M21 12c-2.4 4 -5.4 6 -9 6c-3.6 0 -6.6 -2 -9 -6c2.4 -4 5.4 -6 9 -6c3.6 0 6.6 2 9 6"
                          />
                        </svg>
                        {$showContent[$selectedFile]
                          ? "Hide Code"
                          : "Show Code"}
                      </button>
                      <button
                        class="btn"
                        on:click={() =>
                          copyOutput($obfuscatedContent[$selectedFile])}
                      >
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          width="18"
                          height="18"
                          viewBox="0 0 24 24"
                          fill="none"
                          stroke="currentColor"
                          stroke-width="2"
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          class="icon icon-tabler icons-tabler-outline icon-tabler-copy"
                        >
                          <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                          <path
                            d="M7 7m0 2.667a2.667 2.667 0 0 1 2.667 -2.667h8.666a2.667 2.667 0 0 1 2.667 2.667v8.666a2.667 2.667 0 0 1 -2.667 2.667h-8.666a2.667 2.667 0 0 1 -2.667 -2.667z"
                          />
                          <path
                            d="M4.012 16.737a2.005 2.005 0 0 1 -1.012 -1.737v-10c0 -1.1 .9 -2 2 -2h10c.75 0 1.158 .385 1.5 1"
                          />
                        </svg>
                        Copy Output
                      </button>
                    </div>
                  {:else}
                    <div class="flex ms-auto gap-1 shrink-0 items-start">
                      <button
                        class="btn"
                        on:click={() => obfuscateFile($selectedFile)}
                      >
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          width="18"
                          height="18"
                          viewBox="0 0 24 24"
                          fill="none"
                          stroke="currentColor"
                          stroke-width="2"
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          class="icon icon-tabler icons-tabler-outline icon-tabler-shield-code"
                          ><path
                            stroke="none"
                            d="M0 0h24v24H0z"
                            fill="none"
                          /><path
                            d="M12 21a12 12 0 0 1 -8.5 -15a12 12 0 0 0 8.5 -3a12 12 0 0 0 8.5 3a12 12 0 0 1 -.078 7.024"
                          /><path d="M20 21l2 -2l-2 -2" /><path
                            d="M17 17l-2 2l2 2"
                          /></svg
                        >
                        Obfuscate this file
                      </button>
                    </div>
                  {/if}
                </div>
                <div class="overflow-y-auto grow">
                  {#if !$obfuscatedContent[$selectedFile]}
                    <div class="p-4">
                      <pre
                        class="whitespace-pre-wrap break-words">{$filesContent[
                          $selectedFile
                        ]}
                    </pre>
                    </div>
                  {/if}
                  {#if $showContent[$selectedFile]}
                    <div
                      class="p-4 border-b border-black/15 dark:border-white/15"
                    >
                      <pre
                        class="whitespace-pre-wrap break-words">{$filesContent[
                          $selectedFile
                        ]}
                        </pre>
                    </div>
                  {/if}
                  {#if $obfuscatedContent[$selectedFile]}
                    <pre
                      class="p-4 rounded-lg whitespace-pre-wrap break-words">{$obfuscatedContent[
                        $selectedFile
                      ]}</pre>
                  {/if}
                </div>
              </div>
            {/if}
          </div>
        {/if}
      </div>
    </div>
  </div>
</div>
