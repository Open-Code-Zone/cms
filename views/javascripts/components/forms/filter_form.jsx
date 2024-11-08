import React from 'react'
import { Button } from '../ui/button'
import { Input } from '../ui/input'
import { Popover, PopoverContent, PopoverTrigger } from '../ui/popover'
import { FilterIcon } from 'lucide-react'
import { DatePickerInput, PillInput } from './components'

export default function FilterForm({ config }) {
  const [isOpen, setIsOpen] = React.useState(false)
  const [filters, setFilters] = React.useState({})
  const formRef = React.useRef(null);
  const [url, setUrl] = React.useState("/");

  const handleInputChange = (name, value) => {
    console.log(filters)
    setFilters((prev) => ({ ...prev, [name]: value }))
  }

  React.useEffect(() => {
    const urlEncodedData = new URLSearchParams();
    for (const [key, value] of Object.entries(filters)) {
      if (Array.isArray(value)) {
        // Handle array inputs
        value.forEach((item) => urlEncodedData.append(key, item));
      } else {
        urlEncodedData.append(key, value);
      }
    }

    // Navigate to the same route with query parameters
    const url = `/${config.collection}?${urlEncodedData.toString()}`;
    setUrl(url);
  }, [filters])


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
          <DatePickerInput
            key={field.name}
            value={filters[field.name]}
            onSelect={(value) => handleInputChange(field.name, value)}
          />
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
          <a href={url} rel="noopener noreferrer">
            <Button
              className="bg-blue-500 w-full"
            >
              Apply Filters
            </Button>
          </a>
        </PopoverContent>
      </Popover>
    </div>
  )
}
