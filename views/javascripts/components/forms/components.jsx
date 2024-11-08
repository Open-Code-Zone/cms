import { Popover, PopoverContent, PopoverTrigger } from "../ui/popover";
import { Calendar } from "../ui/calender";
import { format } from "date-fns";
import { CalendarIcon, XIcon } from "lucide-react";
import React from "react";
import { Button } from "../ui/button";

// Helper components
export const DatePickerInput = ({ value, onSelect }) => {
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

export const EditableTitle = ({ value, onChange, disabled }) => {
  const [title, setTitle] = React.useState(value);
  const textareaRef = React.useRef(null);

  const adjustTextareaHeight = () => {
    const textarea = textareaRef.current;
    if (textarea) {
      textarea.style.height = 'auto';
      const newHeight = Math.max(textarea.scrollHeight, 36); // Adjust 36 to match line height if needed
      textarea.style.height = `${newHeight}px`;
    }
  };

  const handleInputChange = (e) => {
    setTitle(e.target.value);
    adjustTextareaHeight();
    onChange(e);
  };

  React.useEffect(() => {
    adjustTextareaHeight();
    window.addEventListener('resize', adjustTextareaHeight);
    return () => window.removeEventListener('resize', adjustTextareaHeight);
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
      disabled={disabled}
    />
  );
};


export const PillInput = ({ value, onChange, placeholder, disabled }) => {
  const [inputValue, setInputValue] = React.useState('');

  const handleKeyDown = (e) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      handleAddValue();
    }
  };

  const handleAddValue = () => {
    if (disabled) return; // Prevent adding pills when disabled
    if (inputValue.trim()) {
      if (value.includes(inputValue.trim())) return;
      onChange([...value, inputValue.trim()]);
      setInputValue('');
    }
  }

  const handleRemove = (index) => {
    if (disabled) return; // Prevent removing pills when disabled
    const newValue = value.filter((_, i) => i !== index);
    onChange(newValue);
  };

  return (
    <div className="flex flex-wrap gap-2 p-2 border rounded-md">
      {value.map?.((item, index) => (
        <span key={index} className="flex items-center px-2 py-1 text-sm bg-blue-200 text-secondary-foreground rounded-md">
          {item}
          {!disabled && (
            <button
              onClick={() => handleRemove(index)}
              className="ml-1 text-secondary-foreground hover:text-primary-foreground hover:bg-red-500 p-0.5 rounded-full focus:outline-none"
              aria-label="Remove pill"
            >
              <XIcon size={12} />
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
        className="flex-grow bg-white text-sm placeholder:text-muted-foreground disabled:cursor-not-allowed disabled:bg-white border-none p-1 focus:outline-none focus:ring-0"
        disabled={disabled} // Disable the input when disable is true
      />
    </div>
  );
};

