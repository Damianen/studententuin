import React from "react";
import { Link } from "react-router-dom";

const Contact = (props) => {
  return (
    <div className="flex min-h-full flex-col justify-center px-6 py-12 lg:px-8 space-y-2">
      <div className="sm:mx-auto sm:w-full sm:max-w-sm">
        <h2 className="text-center text-2xl font-bold leading-9 tracking-tight text-gray-900">Contact</h2>
      </div>
       <div className="mt-10 sm:mx-auto sm:w-full sm:max-w-sm shadow-xl rounded-lg p-10 border border-gray-200">
          <form className="space-y-4" action="#" method="POST">
            <div>
              <label htmlFor="email" className="block text-sm font-medium leading-6 text-gray-900">Jouw Email</label>
                <div className="mt-2">
                  <input id="email" name="email" type="email" required className="block w-full rounded-md border-0 p-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-primary-green sm:text-sm sm:leading-6 focus:outline-none"/>
                </div>
              </div>
            <div>
              <label htmlFor="subject" className="block text-sm font-medium leading-6 text-gray-900">Onderwerp</label>    
              <div className="mt-2">
                <input id="text" name="subject" type="text" required className="block w-full rounded-md border-0 p-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-primary-green sm:text-sm sm:leading-6 focus:outline-none"/>
              </div>
            </div>
            <div>
              <label htmlFor="subject" className="block text-sm font-medium leading-6 text-gray-900">Jouw Bericht</label>    
              <div className="mt-2">
                <textarea id="text" name="subject" type="text" required className="block w-full rounded-md border-0 p-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-primary-green sm:text-sm sm:leading-6 focus:outline-none"> </textarea>
              </div>
            </div>
            <div className="flex justify-end">
              <button type="submit" className="flex justify-center rounded-md bg-primary-green px-3 py-1.5 text-sm font-semibold leading-6 text-white shadow-sm hover:bg-house-green focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-house-green focus:outline-none self-end">Stuur bericht</button>
            </div>
        </form>
      </div>
    </div>
  );
};

export default Contact;