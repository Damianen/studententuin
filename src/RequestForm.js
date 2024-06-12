import React from "react";
import { useLocation } from "react-router-dom";

export default function RequestForm() {
  const location = useLocation();

  // Directly initialize the state from the URL parameters
  const queryParams = new URLSearchParams(location.search);
  const initialPackage = queryParams.get("package") || "";
  const [selectedPackage, setSelectedPackage] = React.useState(initialPackage);

  React.useEffect(() => {
    setInputValue(inputRef.current.value);
    const queryParams = new URLSearchParams(location.search);
    const packageParam = queryParams.get("package");
    if (packageParam) {
      if (packageParam !== selectedPackage) {
        setSelectedPackage(packageParam);
      }
    }
  }, [location.search, selectedPackage]);

  const [formData, setFormData] = React.useState({});
  const [errors, setErrors] = React.useState({});
  const [SuccesMessage, setSuccesMessage] = React.useState({});

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData({
      ...formData,
      [name]: value,
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      setErrors({});
      setSuccesMessage({});

      const response = await fetch("/requestForm", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(formData),
      });

      if (!response.ok) {
        // Assuming the response is a JSON object containing errors
        const errorData = await response.json();
        setErrors(errorData);
      } else {
        const successData = await response.json();
        setSuccesMessage(successData);
        // Handle success (e.g., redirect to another page)
      }
    } catch (error) {
      console.error("Error submitting form:", error);
      setErrors({
        general:
          "Er is een onverwachte fout opgetreden, probeer later opnieuw of neem contact op met de beheerders",
      });
    }
  };

  const handleInputChange = (event) => {
    setInputValue(event.target.value);
  };

  // Combine handlers
  const handleCombinedChange = (e) => {
    handleInputChange(e);
    handleChange(e);
  };

  const handlePackage = (e) => {
    setSelectedPackage(e.target.value);
  };

  // Combine handlers
  const handlePackageAndChange = (e) => {
    handlePackage(e);
    handleChange(e);
  };

  const inputRef = React.useRef(null);

  const [inputValue, setInputValue] = React.useState("");

  return (
    <div className="flex min-h-full flex-col justify-center px-6 py-12 lg:px-8 space-y-2">
      <div className="sm:mx-auto sm:w-full sm:max-w-sm">
        <h2 className="text-center text-2xl font-bold leading-9 tracking-tight text-gray-900">
          Subdomein aanvragen
        </h2>
      </div>
      {errors.general && (
        <div className="text-center text-red-600 font-bold">
          {errors.general}
        </div>
      )}
      {errors.message && (
        <div className="text-center text-red-600 font-bold">
          {errors.message}
        </div>
      )}
      {SuccesMessage.message && (
        <div className="text-center text-green-700 font-bold">
          {SuccesMessage.message}
        </div>
      )}
      <div className="mt-10 sm:mx-auto sm:w-full sm:max-w-sm shadow-xl rounded-lg p-10 border border-gray-200">
        <form className="space-y-4" onSubmit={handleSubmit}>
          <div>
            <label
              htmlFor="subdomainName"
              className="block text-sm font-medium leading-6 text-gray-900"
            >
              Subdomein Naam
            </label>
            <div className="mt-2">
              <input
                id="subdomainName"
                name="subdomainName"
                type="subdomainName"
                required
                className="p-2 block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-primary-green sm:text-sm sm:leading-6 focus:outline-none"
                onChange={handleCombinedChange}
                ref={inputRef}
              />
              <div className="w-full break-words">
                <p>
                  Voorbeeld:{" "}
                  <label className="font-bold">
                    {inputValue}.studententuin.nl
                  </label>
                </p>
              </div>
            </div>
          </div>

          <div>
            <label
              htmlFor="emailAddress"
              className="block text-sm font-medium leading-6 text-gray-900"
            >
              Email
            </label>
            <div className="mt-2">
              <input
                id="emailAddress"
                name="emailAddress"
                type="email"
                required
                onChange={handleChange}
                className="p-2 block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-primary-green sm:text-sm sm:leading-6 focus:outline-none"
              />
            </div>
          </div>

          <div>
            <label
              htmlFor="password"
              className="block text-sm font-medium leading-6 text-gray-900"
            >
              Wachtwoord
            </label>
            <div className="mt-2">
              <input
                id="password"
                name="password"
                type="password"
                required
                onChange={handleChange}
                className="p-2 block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-primary-green sm:text-sm sm:leading-6 focus:outline-none"
              />
            </div>
            <ul>
              <li>Tussen de 8 - 20 karakters</li>
              <li>Een hoofdletter</li>
              <li>Een kleine letter</li>
              <li>Een cijfer</li>
            </ul>
          </div>

          <div>
            <label
              htmlFor="productPackage"
              className="block text-sm font-medium leading-6 text-gray-900"
            >
              Product pakket
            </label>
            <div className="mt-2">
              <select
                name="productPackage"
                className="p-1 block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-primary-green sm:text-sm sm:leading-6 focus:outline-none"
                id="productPackage"
                required
                value={selectedPackage ? JSON.parse(selectedPackage) : ""}
                onChange={handlePackageAndChange}
              >
                <option value="">Kies een pakket</option>
                <option value="free">Gratis</option>
                <option value="basic">Basis</option>
                <option value="premium">Premium</option>
              </select>
            </div>
          </div>

          <div className="flex justify-end">
            <button
              type="submit"
              className="flex justify-center rounded-md bg-primary-green px-3 py-1.5 text-sm font-semibold leading-6 text-white shadow-sm hover:bg-house-green focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-house-green focus:outline-none"
            >
              Aanvraag opsturen
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
