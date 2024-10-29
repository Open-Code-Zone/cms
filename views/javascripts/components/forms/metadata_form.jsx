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

  return (
    <Popover>
      <PopoverTrigger asChild>
        <Button
          variant={"outline"}
          className={`w-[280px] justify-start text-left font-normal ${!date ? "text-muted-foreground" : ""}`}
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
            const formattedDate = `${v.getFullYear()}-${String(v.getMonth() + 1).padStart(2, '0')}-${String(v.getDate()).padStart(2, '0')}`;
            onSelect(v);
          }}
          initialFocus
        />
      </PopoverContent>
    </Popover>
  );
};

const EditableTitle = ({ t, onChange }) => {
  const [title, setTitle] = React.useState(t);
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

const PillInput = ({ value, onChange, placeholder }) => {
  const [inputValue, setInputValue] = React.useState('');

  const handleKeyDown = (e) => {
    if (e.key === 'Enter' && inputValue.trim()) {
      e.preventDefault();
      onChange([...value, inputValue.trim()]);
      setInputValue('');
    }
  };

  const handleRemove = (index) => {
    const newValue = value.filter((_, i) => i !== index);
    onChange(newValue);
  };

  return (
    <div className="flex flex-wrap gap-2 p-2 border rounded-md">
      {value.map((item, index) => (
        <span key={index} className="flex items-center pl-4 pr-2 py-1 text-sm bg-blue-600 text-primary-foreground rounded-full">
          {item}
          <button
            onClick={() => handleRemove(index)}
            className="ml-2 text-primary-foreground hover:text-red-600 focus:outline-none"
            aria-label="Remove pill"
          >
            <X size={16} />
          </button>
        </span>
      ))}
      <input
        type="text"
        value={inputValue}
        onChange={(e) => setInputValue(e.target.value)}
        onKeyDown={handleKeyDown}
        placeholder={placeholder}
        className="flex-grow text-sm placeholder:text-muted-foreground border-none p-1 focus:outline-none focus:ring-0"
      />
    </div>
  );
};

export default function MetaDataForm({ frontMatter, setFrontMatter }) {
  // Initialize config state
  const [config, setConfig] = React.useState(null);
  const [fileNameFormat, setFileNameFormat] = React.useState("");
  const [metadataFields, setMetadataFields] = React.useState([]);
  const [initialFileName, setInitialFileName] = React.useState("");

  // Initialize configuration on component mount
  React.useEffect(() => {
    const fileNameInput = document.getElementById('fileName');
    if (fileNameInput) {
      try {
        const blogConfig = jsyaml.load(fileNameInput.getAttribute("data-blog-config"));
        setConfig(blogConfig);
        setFileNameFormat(blogConfig.file_name_format);
        setInitialFileName(fileNameInput.value);

        // Extract metadata fields from format
        const fields = blogConfig.file_name_format.match(/{(.*?)}/g)?.map(field => field.slice(1, -1)) || [];
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

  const renderInputField = (field) => {
    switch (field.type) {
      case 'string':
        if (field.name === 'description') {
          return (
            <Textarea
              key={field.name}
              value={frontMatter[field.name] || ''}
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
              t={frontMatter[field.name] || ''}
              onChange={(e) => handleInputChange(field.name, e.target.value)}
            />
          );
        }
        return (
          <Input
            key={field.name}
            value={frontMatter[field.name] || ''}
            onChange={(e) => handleInputChange(field.name, e.target.value)}
            placeholder={field.description}
            className="w-full"
          />
        );
      case 'datetime':
        return (
          <DatePickerInput
            key={field.name}
            value={frontMatter[field.name]}
            onSelect={(value) => handleInputChange(field.name, value)}
          />
        );
      case 'array':
        return (
          <PillInput
            key={field.name}
            value={frontMatter[field.name] || []}
            onChange={(value) => handleInputChange(field.name, value)}
            placeholder={field.description}
          />
        );
      case 'file':
        return (
          <div key={field.name} className="relative">
            <Input
              value={frontMatter[field.name] || ''}
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
  if (!config) {
    return null;
  }

  return (
    <form className="max-w-2xl mx-auto space-y-6 p-6 bg-background">
      {config.metadata_schema.map((field) => (
        <div key={field.name}>
          {renderInputField(field)}
        </div>
      ))}
    </form>
  );
}
