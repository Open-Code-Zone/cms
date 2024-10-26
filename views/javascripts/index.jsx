import React from 'react';
import "../css/prosemirror.css";
import { createRoot } from 'react-dom/client';
import Editor from "./components/editor/advanced-editor";
import MetaDataForm from './components/metadata-form/form';
import jsyaml from 'js-yaml';

const editorContainer = document.getElementById('editor');
const root = createRoot(editorContainer);
const textarea = document.getElementById('contentInput');
const initialValue = textarea.value;

const id = textarea.getAttribute("data-id");
const submitButton = document.getElementById('submit-button');

const loadContentFromLocalStorage = (id) => {
  const savedContent = window.localStorage.getItem(id);
  return savedContent || textarea.value || '';
};

const parseMarkdownContent = (fileContent) => {
  if (fileContent === "") return { frontMatter: "", markdownContent: "" }

  const [_, fM, mC] = fileContent.split('---');
  try {
    return { frontMatter: jsyaml.load(fM), markdownContent: mC };
  } catch (error) {
    console.error("Error parsing YAML:", error);
    return { frontMatter: {}, markdownContent: mC };
  }
};

const saveContentToLocalStorage = (id, fileContent) => {
  window.localStorage.setItem(id, fileContent);
};

const generateFileContent = (frontMatter, markdownContent) => {
  const yaml = jsyaml.dump(frontMatter);
  // remove the first newline character from the markdown content
  // we don't want to replace all newlines only the starting one
  try {
    markdownContent = markdownContent.replace(/^\n/, '');
  } catch (err) {
    console.log(err)
  }

  return `---\n${yaml}---\n${markdownContent}`;
};

const initialContent = parseMarkdownContent(loadContentFromLocalStorage(id));

const NovelEditor = () => {
  const [markdownContent, setMarkdownContent] = React.useState(initialContent.markdownContent);
  const [frontMatter, setFrontMatter] = React.useState(initialContent.frontMatter);

  React.useEffect(() => {
    const fileContent = generateFileContent(frontMatter, markdownContent);
    textarea.value = fileContent;
    saveContentToLocalStorage(id, fileContent);

    if (fileContent !== initialValue) {
      submitButton.disabled = false;
    } else {
      submitButton.disabled = true;
    }
  }, [frontMatter, markdownContent]);

  return (
    <>
      <MetaDataForm frontMatter={frontMatter} setFrontMatter={setFrontMatter} />
      <Editor markdownContent={markdownContent} setMarkdownContent={setMarkdownContent} />
    </>
  );
};

root.render(<NovelEditor />);
