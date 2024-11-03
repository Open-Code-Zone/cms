import * as React from "react";
import { Input } from '../ui/input';
import { Textarea } from "../ui/textarea";
import { Button } from "../ui/button";
import { X, CalendarIcon, Upload } from "lucide-react";
import jsyaml from 'js-yaml';
import { Popover, PopoverContent, PopoverTrigger } from "../ui/popover";
import { Calendar } from "../ui/calender";
import { format } from "date-fns";

// Helper components
const DatePickerInput = ({ value, onSelect }) => {
  const [date, setDate] = React.useState(value ? new Date(value) : null);
  const [open, setOpen] = React.useState(false);

  return (
    <Popover open={open}>
      <PopoverTrigger asChild>
        <Button
          variant={"outline"}
          className={`w-[280px] justify-start text-left font-normal ${!date ? "text-muted-foreground" : ""}`}
          onClick={() => setOpen(!open)}
        >
          <CalendarIcon className="mr-2 h-4 w-4" />
          {date ? format(date, "PPP") : <span>Pick a date</span>}
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-auto p-0">
        <Calendar
          mode="single"
          selected={date}
          onSelect={(v) => {
            setDate(v);
            //const formattedDate = `${v.getFullYear()}-${String(v.getMonth() + 1).padStart(2, '0')}-${String(v.getDate()).padStart(2, '0')}`;
            onSelect(v);
            setOpen(false);
          }}
          initialFocus
        />
      </PopoverContent>
    </Popover >
  );
};

const EditableTitle = ({ value, onChange }) => {
  const [title, setTitle] = React.useState(value);
  const textareaRef = React.useRef(null);

  const adjustTextareaHeight = () => {
    const textarea = textareaRef.current;
    if (textarea) {
      textarea.style.height = 'auto';
      textarea.style.height = `${textarea.scrollHeight}px`;
    }
  };

  const handleInputChange = (e) => {
    setTitle(e.target.value);
    adjustTextareaHeight();
    onChange(e);
  };

  React.useEffect(() => {
    adjustTextareaHeight();
    const handleResize = () => adjustTextareaHeight();
    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, []);

  return (
    <textarea
      ref={textareaRef}
      value={title}
      onChange={handleInputChange}
      onKeyDown={(e) => {
        if (e.key === 'Enter') {
          e.preventDefault();
        }
      }}
      placeholder="Title"
      className="text-4xl font-bold border-none focus:outline-none p-0 w-full resize-none overflow-hidden bg-transparent"
      style={{ lineHeight: '1.2' }}
    />
  );
}

const PillInput = ({ value, onChange, placeholder, disable }) => {
  const [inputValue, setInputValue] = React.useState('');

  const handleKeyDown = (e) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      handleAddValue();
    }
  };

  const handleAddValue = () => {
    if (disable) return; // Prevent adding pills when disabled
    if (inputValue.trim()) {
      if (value.includes(inputValue.trim())) return;
      onChange([...value, inputValue.trim()]);
      setInputValue('');
    }
  }

  const handleRemove = (index) => {
    if (disable) return; // Prevent removing pills when disabled
    const newValue = value.filter((_, i) => i !== index);
    onChange(newValue);
  };

  return (
    <div className="flex flex-wrap gap-2 p-2 border rounded-md">
      {value.map?.((item, index) => (
        <span key={index} className="flex items-center px-3 py-1 text-sm bg-blue-600 text-primary-foreground rounded-md">
          {item}
          {!disable && (
            <button
              onClick={() => handleRemove(index)}
              className="ml-2 text-primary-foreground hover:text-red-600 focus:outline-none"
              aria-label="Remove pill"
            >
              <X size={16} />
            </button>
          )}
        </span>
      ))}
      <input
        type="text"
        value={inputValue}
        onChange={(e) => setInputValue(e.target.value)}
        onKeyDown={handleKeyDown}
        onBlur={handleAddValue}
        placeholder={placeholder}
        className="flex-grow text-sm placeholder:text-muted-foreground border-none p-1 focus:outline-none focus:ring-0"
        disabled={disable} // Disable the input when disable is true
      />
    </div>
  );
};

export default function MetaDataForm({ frontMatter, setFrontMatter }) {
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
        setCollectionConfig(jsyaml.load(fileNameInput.getAttribute("data-collection-config")));
        setCollectionPermission(jsyaml.load(fileNameInput.getAttribute("data-user-config")));
        //setFileNameFormat(collectionConfig.file_name_format);
        setInitialFileName(fileNameInput.value);

        // Extract metadata fields from format
        //const fields = collectionConfig.file_name_format.match(/{(.*?)}/g)?.map(field => field.slice(1, -1)) || [];
        //setMetadataFields(fields);
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

  const renderInputField = (field) => {
    switch (field.type) {
      case 'string':
        if (field.name === 'description') {
          return (
            <Textarea
              key={field.name}
              value={frontMatter[field.name] || collectionPermission.default_metadata[field.name] || ''}
              onChange={(e) => handleInputChange(field.name, e.target.value)}
              placeholder={field.description}
              className="w-full"
              rows={4}
            />
          );
        } else if (field.name === 'title') {
          return (
            <EditableTitle
              key={field.name}
              value={frontMatter[field.name] || collectionPermission.default_metadata[field.name] || ''}
              onChange={(e) => handleInputChange(field.name, e.target.value)}
            />
          );
        }
        return (
          <Input
            key={field.name}
            value={frontMatter[field.name] || collectionPermission.default_metadata[field.name] || ''}
            onChange={(e) => handleInputChange(field.name, e.target.value)}
            placeholder={field.description}
            className="w-full"
          />
        );
      case 'datetime':
        return (
          <DatePickerInput
            key={field.name}
            value={frontMatter[field.name] || collectionPermission.default_metadata[field.name]}
            onSelect={(value) => handleInputChange(field.name, value)}
          />
        );
      case 'array':
        return (
          <PillInput
            key={field.name}
            value={frontMatter[field.name] || collectionPermission.default_metadata[field.name]?.value || []}
            disable={collectionPermission.default_metadata[field.name]?.strict || false}
            onChange={(value) => handleInputChange(field.name, value)}
            placeholder={field.description}
          />
        );
      case 'file':
        return (
          <div key={field.name} className="relative">
            <Input
              value={frontMatter[field.name] || collectionPermission.default_metadata[field.name] || ''}
              onChange={(e) => handleInputChange(field.name, e.target.value)}
              placeholder={field.description}
              className="py-5 px-2.5 pr-24"
            />
            <Button
              type="button"
              variant="outline"
              size="sm"
              className="absolute right-1 top-1 flex items-center"
              onClick={() => handleImageUpload(field.name)}
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
      {collectionConfig.metadata_schema.map((field) => (
        <div key={field.name}>
          {renderInputField(field)}
        </div>
      ))}
    </form>
  );
}
