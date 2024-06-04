import React from "react";
import { Link } from "react-router-dom";

const Login = (props) => {
  return (
    <div class="flex justify-center items-center h-screen ">
      <div class="justify-center items-center">
        <h1 class="mt-10 text-center text-2xl font-bold leading-9 tracking-tight text-gray-900 p-5">
          Login
        </h1>
        <div class="lg:p-28 md:p-36 sm:20 p-4 w-full shadow-md border border-gray-200 rounded-lg">
          <form action="/login" method="POST">
            {/* <!-- Email Input --> */}
            <div class="mb-4">
              <label
                for="email"
                class="block text-sm font-medium leading-6 text-gray-900"
              >
                Email
              </label>
              <input
                type="email"
                id="email"
                name="email"
                class="w-full border border-gray-300 rounded-md py-2 px-3 focus:outline-none focus:border-green-400"
                autocomplete="off"
              ></input>
            </div>
            {/* <!-- Password Input --> */}
            <div class="mb-6">
              <label
                for="password"
                class="block text-sm font-medium leading-6 text-gray-900"
              >
                Password
              </label>
              <input
                type="password"
                id="password"
                name="password"
                class="w-full border border-gray-300 rounded-md py-2 px-3 focus:outline-none focus:border-green-400"
              ></input>
            </div>
            {/* <!-- Remember me and Forgot password Link --> */}
            <div className="flex items-center">
              <input
                type="checkbox"
                id="remember"
                name="remember"
                className="text-green-400"
              ></input>
              <label htmlFor="remember" className="text-gray-600 ml-2">
                Onthoud mij
              </label>
            </div>
            <div className="text-primary-green">
              <Link to="wachtwoord" className="hover:underline">
                Wachtwoord vergeten?
              </Link>
            </div>
            <div class="mb-6 text-primary-green">
              <a href="#" class="hover:underline text-primary-green h-100%">
                Sign up Here
              </a>
            </div>
            {/* <!-- Sign up  Link --> */}
            <div className="flex justify-between items-center">
              <button
                type="submit"
                class=" shadow-md rounded-full bg-gray-500 px-5 py-2 text-sm font-semibold leading-6 text-white shadow-sm hover:bg-gray-600 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-gray-600"
              >
                Back{" "}
              </button>
              {/* <!-- Login Button --> */}
              <button
                type="submit"
                class=" shadow-md rounded-full bg-primary-green px-5 py-2 text-sm font-semibold leading-6 text-white shadow-sm hover:bg-house-green focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-green-600"
              >
                Login
              </button>
            </div>
        </form>
      </div>
    </div>
  </div>
  );
};

export default Login;
