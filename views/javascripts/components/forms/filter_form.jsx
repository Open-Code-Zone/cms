import React from 'react'
import { Button } from '../ui/button'
import { Input } from '../ui/input'
import { Popover, PopoverContent, PopoverTrigger } from '../ui/popover'
import { Calendar } from '../ui/calender'
import { CalendarIcon, FilterIcon, X } from 'lucide-react'
import { format } from 'date-fns'

const PillInput = ({ value, onChange, placeholder, name }) => {
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
        <React.Fragment key={index}>
          <span className="flex items-center pl-4 pr-2 py-1 text-sm font-bold bg-green-500 text-primary-foreground rounded-full">
            {item}
            <button
              onClick={() => handleRemove(index)}
              className="ml-2 text-primary-foreground hover:text-black hover:bg-white rounded-full focus:outline-none"
              aria-label="Remove pill"
            >
              <X size={16} />
            </button>
          </span>
          <input
            type="hidden"
            name={name}
            value={item}
          />
        </React.Fragment>
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

export default function FilterForm({ config }) {
  const [isOpen, setIsOpen] = React.useState(false)
  const [filters, setFilters] = React.useState({})
  const formRef = React.useRef(null);

  const handleInputChange = (name, value) => {
    setFilters((prev) => ({ ...prev, [name]: value }))
  }

  const handleSubmit = (e) => {
    e.preventDefault();
    const form = formRef.current;
    // convert form data to url encoded string
    const formData = new FormData(form);
    const urlEncodedData = new URLSearchParams();
    for (const pair of formData) {
      console.log(pair)
      urlEncodedData.append(pair[0], pair[1]);
    }

    if (window.htmx && form) {
      window.htmx.ajax('POST', `/${config.collection}/filter?${urlEncodedData.toString()}`, {
        target: '#results',
        swap: 'outerHTML',
        values: new FormData(form),
        headers: {
          'HX-Request': 'true'
        }
      });
    }
  }

  const renderFilterField = (field) => {
    switch (field.type) {
      case 'string':
        return (
          <Input
            key={field.name}
            name={field.name}
            placeholder={field.description}
            value={filters[field.name] || ''}
            onChange={(e) => handleInputChange(field.name, e.target.value)}
          />
        );
      case 'datetime':
        return (
          <div key={field.name}>
            <Popover>
              <PopoverTrigger asChild>
                <Button variant="outline">
                  <CalendarIcon className="mr-2 h-4 w-4" />
                  {filters[field.name] ? (
                    format(new Date(filters[field.name]), 'PPP')
                  ) : (
                    <span>Pick a date</span>
                  )}
                </Button>
              </PopoverTrigger>
              <PopoverContent className="w-auto p-0">
                <Calendar
                  mode="single"
                  selected={filters[field.name] ? new Date(filters[field.name]) : undefined}
                  onSelect={(date) => handleInputChange(field.name, date ? date.toISOString() : '')}
                  initialFocus
                />
              </PopoverContent>
            </Popover>
            <input
              type="hidden"
              name={field.name}
              value={filters[field.name] || ''}
            />
          </div>
        );
      case 'array':
        return (
          <PillInput
            key={field.name}
            name={field.name}
            value={filters[field.name] || []}
            onChange={(value) => handleInputChange(field.name, value)}
            placeholder={field.description}
          />
        );
      default:
        return null;
    }
  };

  return (
    <div>
      <Popover open={isOpen} onOpenChange={setIsOpen}>
        <PopoverTrigger asChild>
          <Button variant="outline">
            <FilterIcon className="mr-2 h-4 w-4" />
            Filter
          </Button>
        </PopoverTrigger>
        <PopoverContent className="w-80">
          <form
            ref={formRef}
            className="space-y-4"
            onSubmit={handleSubmit}
          >
            {config.metadata_schema.map((field) => (
              field.filter ? (
                <div key={field.name} className="space-y-2">
                  <label htmlFor={field.name} className="text-sm font-medium">
                    {field.name}
                  </label>
                  {renderFilterField(field)}
                </div>
              ) : null
            ))}
            <Button
              type="submit"
              className="w-full"
            >
              Apply Filters
            </Button>
          </form>
        </PopoverContent>
      </Popover>
    </div>
  )
}
