import React, {useState } from 'react';
import "../css/prosemirror.css";
import { createRoot } from 'react-dom/client';
import Editor from "./components/editor/advanced-editor";

const editorContainer = document.getElementById('editor');
const root = createRoot(editorContainer);

const NovelEditor = () => {
  const [value, setValue] = useState({});
  console.log(value);
  return (
    <Editor initialValue={value} onChange={setValue} />
  );
};

root.render(<NovelEditor />);
