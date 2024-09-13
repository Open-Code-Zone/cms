import * as React from "react";
import { useEditor } from "novel";
import { Check, LinkIcon, Trash } from "lucide-react";
import {
  useEffect,
  useRef,
} from "react";
import { Button } from "../../ui/button";
import { Popover, PopoverContent, PopoverTrigger } from "../../ui/popover";
import { cn } from "../../../utils";

export function isValidUrl(url) {
  try {
    new URL(url);
    return true;
  } catch (e) {
    return false;
  }
}
export function getUrlFromString(str) {
  if (isValidUrl(str)) return str;
  try {
    if (str.includes(".") && !str.includes(" ")) {
      return new URL(`https://${str}`).toString();
    }
  } catch (e) {
    return null;
  }
}

export const LinkSelector = ({ open, onOpenChange }) => {
  const inputRef = useRef(null);
  const { editor } = useEditor();

  useEffect(() => {
    if (open && inputRef.current) {
      inputRef.current.focus();
    }
  }, [open]);

  if (!editor) return null;
  return (
    <Popover modal={true} open={open} onOpenChange={onOpenChange}>
      <PopoverTrigger asChild>
        <Button
          size="sm"
          variant="ghost"
          className="gap-2 rounded-none border-none"
        >
          <p className="text-base">
            <LinkIcon size={12} />
          </p>
          <p
            className={cn("underline decoration-stone-400 underline-offset-4", {
              "text-blue-500": editor.isActive("link"),
            })}
          >
            Link
          </p>
        </Button>
      </PopoverTrigger>
      <PopoverContent align="start" className="w-60 p-0" sideOffset={10}>
        <form
          onSubmit={(e) => {
            const target = e.currentTarget;
            e.preventDefault();
            const input = target[0];
            const url = getUrlFromString(input.value);
            if (url) {
              editor.chain().focus().setLink({ href: url }).run();
              onOpenChange(false);
            }
          }}
          className="flex  p-1"
        >
          <input
            ref={inputRef}
            type="text"
            placeholder="Paste a link"
            className="flex-1 bg-background p-1 text-sm outline-none"
            defaultValue={editor.getAttributes("link").href || ""}
          />
          {editor.getAttributes("link").href ? (
            <Button
              size="icon"
              variant="outline"
              type="button"
              className="flex h-8 items-center rounded-sm p-1 text-red-600 transition-all hover:bg-red-100 dark:hover:bg-red-800"
              onClick={() => {
                editor.chain().focus().unsetLink().run();
                onOpenChange(false);
              }}
            >
              <Trash size={16} />
            </Button>
          ) : (
            <Button size="icon" className="h-8">
              <Check size={16} />
            </Button>
          )}
        </form>
      </PopoverContent>
    </Popover>
  );
};
