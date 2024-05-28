import React from "react";

const Login = (props) => {
  return (
    <div className="flex min-h-full flex-col justify-center px-6 py-12 lg:px-8 space-y-2">
      <div className="sm:mx-auto sm:w-full sm:max-w-sm">
        <h2 className="mt-10 text-center text-2xl font-bold leading-9 tracking-tight text-gray-900">Login</h2>
      </div>
       <div className="mt-10 sm:mx-auto sm:w-full sm:max-w-sm shadow-xl rounded-lg p-10 border border-gray-200">
          <form className="space-y-4" action="#" method="POST">
            <div>
              <label htmlFor="email" className="block text-sm font-medium leading-6 text-gray-900">Email</label>
                <div className="mt-2">
                  <input id="email" name="email" type="email" required className="block w-full rounded-md border-0 p-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-primary-green sm:text-sm sm:leading-6 focus:outline-none"/>
                </div>
              </div>
            <div>
              <label htmlFor="password" className="block text-sm font-medium leading-6 text-gray-900">Wachtwoord</label>    
              <div className="mt-2">
                <input id="password" name="password" type="password" required className="block w-full rounded-md border-0 p-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-primary-green sm:text-sm sm:leading-6 focus:outline-none"/>
              </div>
            </div>
            <div className="flex items-center">
              <input
                type="checkbox"
                id="remember"
                name="remember"
                className="text-green-400"
              ></input>
              <label htmlFor="remember" className="text-gray-600 ml-2">
                Remember Me
              </label>
            </div>
            <div className="text-primary-green">
              <a href="#" className="hover:underline">
                Forgot Password?
              </a>
            </div>
            <div className="flex justify-between items-center">
              <a href="#" className="hover:underline text-primary-green h-100% self-start">
                Sign up Here
              </a>
              <button type="submit" className="flex justify-center rounded-md bg-primary-green px-3 py-1.5 text-sm font-semibold leading-6 text-white shadow-sm hover:bg-house-green focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-house-green focus:outline-none self-end">Login</button>
            </div>
        </form>
      </div>
    </div>
  );
};

export default Login;
