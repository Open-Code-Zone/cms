import React from 'react';
import "../css/prosemirror.css";
import { createRoot } from 'react-dom/client';
import Editor from "./components/editor/advanced-editor";
import jsyaml from 'js-yaml';
import MetaDataForm from './components/forms/metadata_form';
import FilterForm from './components/forms/filter_form';

// Editor related functions
const loadContentFromLocalStorage = (id) => {
  const savedContent = window.localStorage.getItem(id);
  const textarea = document.getElementById('contentInput');
  return savedContent || (textarea ? textarea.value : '') || '';
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
  try {
    markdownContent = markdownContent.replace(/^\n/, '');
  } catch (err) {
    console.log(err)
  }
  return `---\n${yaml}---\n${markdownContent}`;
};

// Novel Editor Component
const NovelEditor = ({ id, initialValue }) => {
  const initialContent = parseMarkdownContent(loadContentFromLocalStorage(id));
  const [markdownContent, setMarkdownContent] = React.useState(initialContent.markdownContent);
  const [frontMatter, setFrontMatter] = React.useState(initialContent.frontMatter);
  const submitButton = document.getElementById('submit-button');
  const textarea = document.getElementById('contentInput');

  React.useEffect(() => {
    if (!textarea || !submitButton) return;

    const fileContent = generateFileContent(frontMatter, markdownContent);
    textarea.value = fileContent;
    saveContentToLocalStorage(id, fileContent);

    submitButton.disabled = fileContent === initialValue;
  }, [frontMatter, markdownContent, id, initialValue]);

  return (
    <>
      <MetaDataForm frontMatter={frontMatter} setFrontMatter={setFrontMatter} />
      <Editor markdownContent={markdownContent} setMarkdownContent={setMarkdownContent} />
    </>
  );
};

// Initialize components based on presence of DOM elements
const initializeComponents = () => {
  // Try to initialize editor
  const editorContainer = document.getElementById('editor');
  if (editorContainer) {
    const textarea = document.getElementById('contentInput');
    if (textarea) {
      const editorRoot = createRoot(editorContainer);
      const initialValue = textarea.value || "";
      const id = textarea.getAttribute("data-id");

      editorRoot.render(
        <NovelEditor id={id} initialValue={initialValue} />
      );
    }
  }

  // Try to initialize filter form
  const filterForm = document.getElementById('filterFormRoot');
  if (filterForm) {
    try {
      const filterFormRoot = createRoot(filterForm);
      const yamlConfig = jsyaml.load(filterForm.getAttribute('data-config'));
      filterFormRoot.render(<FilterForm config={yamlConfig} />);
    } catch (error) {
      console.error('Error initializing filter form:', error);
    }
  }
};

// Initialize the application
document.addEventListener('DOMContentLoaded', initializeComponents);
