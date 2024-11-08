import * as React from "react";
import { Input } from '../ui/input';
import { Textarea } from "../ui/textarea";
import { Button } from "../ui/button";
import { Upload } from "lucide-react";
import jsyaml from 'js-yaml';
import { DatePickerInput, EditableTitle, PillInput } from "./components";

export default function MetaDataForm({ mode, frontMatter, setFrontMatter }) {
  // Initialize config state
  const [collectionConfig, setCollectionConfig] = React.useState(null);
  const [collectionPermission, setCollectionPermission] = React.useState(null);
  const [fileNameFormat, setFileNameFormat] = React.useState("");
  const [metadataFields, setMetadataFields] = React.useState([]);
  const [initialFileName, setInitialFileName] = React.useState("");

  // Initialize configuration on component mount
  React.useEffect(() => {
    const fileNameInput = document.getElementById('fileName');
    if (fileNameInput) {
      try {
        const collectionConfig = jsyaml.load(fileNameInput.getAttribute("data-collection-config"));
        const collectionPermission = jsyaml.load(fileNameInput.getAttribute("data-user-config"));
        setCollectionConfig(collectionConfig);
        setCollectionPermission(collectionPermission);
        setFileNameFormat(collectionConfig.file_name_format);
        setInitialFileName(fileNameInput.value);

        // Extract metadata fields from format
        const fields = collectionConfig.file_name_format?.match(/{(.*?)}/g)?.map(field => field.slice(1, -1)) || [];
        setMetadataFields(fields);
      } catch (error) {
        console.error('Error loading blog configuration:', error);
      }
    }
  }, []);

  // Handle filename updates
  React.useEffect(() => {
    if (!initialFileName || initialFileName !== "new-draft.md" || !metadataFields.length) return;

    const fileNameInput = document.getElementById('fileName');
    if (fileNameInput) {
      const newFileName = metadataFields
        .map((field) => (frontMatter[field] || '').toString().replace(/\s+/g, '-'))
        .join('-') + '.md';
      fileNameInput.value = newFileName;
    }
  }, [frontMatter, initialFileName, metadataFields]);

  const handleInputChange = (name, value) => {
    setFrontMatter({
      ...frontMatter,
      [name]: value
    });
  };

  const renderInputField = (field, index) => {
    switch (field.type) {
      case 'string':
        if (field.type === 'text') {
          return (
            <Textarea
              key={field.name}
              value={frontMatter[field.name] || collectionPermission.default_metadata?.[field.name] || ''}
              onChange={(e) => handleInputChange(field.name, e.target.value)}
              placeholder={field.description}
              className="w-full"
              rows={4}
              disabled={mode != "write"}
            />
          );
        } else if (field.type === 'string' && index == 0) {
          return (
            <EditableTitle
              key={field.name}
              value={frontMatter[field.name] || collectionPermission.default_metadata?.[field.name] || ''}
              onChange={(e) => handleInputChange(field.name, e.target.value)}
              disabled={mode != "write"}
            />
          );
        }
        return (
          <Input
            key={field.name}
            value={frontMatter[field.name] || collectionPermission.default_metadata?.[field.name] || ''}
            onChange={(e) => handleInputChange(field.name, e.target.value)}
            placeholder={field.description}
            className="w-full"
            disabled={mode != "write"}
          />
        );
      case 'datetime':
        // to do to implement date picker as reader mode
        return (
          <DatePickerInput
            key={field.name}
            value={frontMatter[field.name] || collectionPermission.default_metadata?.[field.name]}
            onSelect={(value) => handleInputChange(field.name, value)}
          />
        );
      case 'array':
        return (
          <PillInput
            key={field.name}
            value={frontMatter[field.name] || collectionPermission.default_metadata?.[field.name]?.value || []}
            disabled={collectionPermission.default_metadata[field.name]?.strict || mode || false}
            onChange={(value) => handleInputChange(field.name, value)}
            placeholder={field.description}
          />
        );
      case 'file':
        return (
          <div key={field.name} className="relative">
            <Input
              value={frontMatter[field.name] || collectionPermission.default_metadata?.[field.name] || ''}
              onChange={(e) => handleInputChange(field.name, e.target.value)}
              placeholder={field.description}
              className="py-5 px-2.5 pr-24"
              disabled={mode != "write"}
            />
            <Button
              type="button"
              variant="outline"
              size="sm"
              className="absolute right-1 top-1 flex items-center"
              onClick={() => handleImageUpload(field.name)}
              disabled={mode != "write"}
            >
              <Upload className="w-4 h-4 mr-2" />
              Upload
            </Button>
          </div>
        );
      default:
        return null;
    }
  };

  const handleImageUpload = (fieldName) => {
    console.log('Image upload triggered for', fieldName);
  };

  // If configuration hasn't loaded yet, return null or a loading state
  if (!collectionConfig) {
    return null;
  }

  return (
    <form className="max-w-2xl mx-auto space-y-6 p-6 bg-background">
      {collectionConfig.metadata_schema.map((field, index) => (
        <div key={field.name}>
          {renderInputField(field, index)}
        </div>
      ))}
    </form>
  );
}
