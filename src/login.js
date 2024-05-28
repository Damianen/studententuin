import React from "react";
// import { useState } from "react";

const Login = (props) => {
  return (
    <div class="flex justify-center items-center h-screen ">
      <div class="justify-center items-center">
        <h1 class="mt-10 text-center text-2xl font-bold leading-9 tracking-tight text-gray-900 p-5">
          Login
        </h1>
        <div class="lg:p-28 md:p-36 sm:20 p-4 w-full shadow-md border border-gray-200 rounded-lg">
          <form action="#" method="POST">
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
            <div class="mb-4">
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
                autocomplete="off"
              ></input>
            </div>
            {/* <!-- Remember Me Checkbox --> */}
            <div class="mb-4 flex items-center">
              <input
                type="checkbox"
                id="remember"
                name="remember"
                class="text-green-400"
              ></input>
              <label for="remember" class="text-gray-600 ml-2">
                Remember Me
              </label>
            </div>
            {/* <!-- Forgot Password Link --> */}
            <div class="mb-6 text-primary-green">
              <a href="#" class="hover:underline">
                Forgot Password?
              </a>
            </div>
            {/* <!-- Sign up  Link --> */}
            <div class="mt-6 text-primary-green">
              <a href="#" class="hover:underline">
                Sign up Here
              </a>
            </div>
            {/* <!-- Login Button --> */}
            <button
              type="submit"
              class="float-right shadow-md rounded-full bg-primary-green px-5 py-2 text-sm font-semibold leading-6 text-white shadow-sm hover:bg-house-green focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-green-600"
            >
              Login
            </button>
          </form>
        </div>
      </div>
    </div>
  );
};

export default Login;
