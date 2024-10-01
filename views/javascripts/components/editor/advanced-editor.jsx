import React, { useEffect, useState } from "react";
import {
  EditorRoot,
  EditorCommand,
  EditorCommandItem,
  EditorCommandEmpty,
  EditorContent,
  EditorCommandList,
  EditorBubble,
} from "novel";
import { ImageResizer, handleCommandNavigation } from "novel/extensions";
import { defaultExtensions } from "./extensions";
import { NodeSelector } from "./selectors/node-selector";
import { LinkSelector } from "./selectors/link-selector";

import { useDebouncedCallback } from "use-debounce";

import { TextButtons } from "./selectors/text-buttons";
import { slashCommand, suggestionItems } from "./slash-command";
import { handleImageDrop, handleImagePaste } from "novel/plugins";
import { uploadFn } from "./image-upload";
import { Separator } from "../ui/separator";

const extensions = [...defaultExtensions, slashCommand];
//const hljs = require('highlight.js');

const Editor = () => {
  const [openNode, setOpenNode] = useState(false);
  const [openLink, setOpenLink] = useState(false);
  const [saveStatus, setSaveStatus] = useState("Draft Saved");
  const [charsCount, setCharsCount] = useState();

  //const highlightCodeblocks = (content) => {
  //  const doc = new DOMParser().parseFromString(content, 'text/html');
  //  doc.querySelectorAll('pre code').forEach((el) => {
  //    // @ts-ignore
  //    // https://highlightjs.readthedocs.io/en/latest/api.html?highlight=highlightElement#highlightelement
  //    hljs.highlightElement(el);
  //  });
  //  return new XMLSerializer().serializeToString(doc);
  //};

  const debouncedUpdates = useDebouncedCallback(async (editor) => {
    setCharsCount(editor.storage.characterCount.words());
    const id = document.getElementById("contentInput").getAttribute("data-id");
    window.localStorage.setItem(id, editor.storage.markdown.getMarkdown());
    document.getElementById("contentInput").value = editor.storage.markdown.getMarkdown();
    setSaveStatus("Draft Saved");
  }, 500);

  const setEditorMarkdownContent = (async (editor) => {
    const id = document.getElementById("contentInput").getAttribute("data-id");
    const markdown = window.localStorage.getItem(id);
    if (markdown) {
      editor.commands.setContent(markdown)
    } else {
      editor.commands.setContent(document.getElementById("contentInput").value)
    }
  })

  return (
    <div className="relative w-full">
      <div className="flex absolute right-5 top-5 z-10 mb-5 gap-2">
        <div className="rounded-lg bg-accent px-2 py-1 text-sm text-muted-foreground">{saveStatus}</div>
        <div className={charsCount ? "rounded-lg bg-accent px-2 py-1 text-sm text-muted-foreground" : "hidden"}>
          {charsCount} Words
        </div>
      </div>

      <EditorRoot>
        <EditorContent
          className="relative min-h-[500px] p-10 w-full border bg-background sm:mb-[calc(20vh)] sm:rounded-lg sm:border"
          extensions={extensions}
          editorProps={{
            handleDOMEvents: {
              keydown: (_view, event) => handleCommandNavigation(event),
            },
            handlePaste: (view, event) => handleImagePaste(view, event, uploadFn),
            handleDrop: (view, event, _slice, moved) => handleImageDrop(view, event, moved, uploadFn),
            attributes: {
              class:
                "prose prose-lg dark:prose-invert prose-headings:font-title font-default focus:outline-none max-w-full",
            },
          }}
          onUpdate={({ editor }) => {
            debouncedUpdates(editor);
            setSaveStatus("Unsaved");
          }}
          onCreate={(editor) => {
            setEditorMarkdownContent(editor.editor);
          }}
          slotAfter={<ImageResizer />}
        >
          <EditorCommand className="z-50 h-auto max-h-[330px] overflow-y-auto rounded-md border border-muted bg-background px-1 py-2 shadow-md transition-all">
            <EditorCommandEmpty className="px-2 text-muted-foreground">No results</EditorCommandEmpty>
            <EditorCommandList>
              {suggestionItems.map((item) => (
                <EditorCommandItem
                  value={item.title}
                  onCommand={(val) => item.command?.(val)}
                  className={`flex cursor-pointer w-full items-center space-x-2 rounded-md px-2 py-1 text-left text-sm hover:bg-accent aria-selected:bg-accent `}
                  key={item.title}
                >
                  <div className="flex h-10 w-10 items-center justify-center rounded-md border border-muted bg-background">
                    {item.icon}
                  </div>
                  <div>
                    <p className="font-medium">{item.title}</p>
                    <p className="text-xs text-muted-foreground">
                      {item.description}
                    </p>
                  </div>
                </EditorCommandItem>
              ))}
            </EditorCommandList>
          </EditorCommand>

          <EditorBubble
            tippyOptions={{
              placement: "top",
            }}
            className="flex bg-white w-fit max-w-[90vw] overflow-hidden rounded-md border shadow-xl"
          >
            <Separator orientation="vertical" />
            <NodeSelector open={openNode} onOpenChange={setOpenNode} />
            <Separator orientation="vertical" />

            <LinkSelector open={openLink} onOpenChange={setOpenLink} />
            <Separator orientation="vertical" />
            <TextButtons />
            <Separator orientation="vertical" />
          </EditorBubble>
        </EditorContent>
      </EditorRoot>
    </div >
  );
};

export default Editor;
