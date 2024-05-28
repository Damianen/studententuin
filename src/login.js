import React from "react";
// import { useState } from "react";

const Login = (props) => {
  return (
    <div class="bg-gray-100 flex justify-center items-center h-screen">
      {/* <!-- Left: Image --> */}
      {/* <div class="w-1/2 h-screen hidden lg:block">
        <img
          src="https://placehold.co/800x/667fff/ffffff.png?text=Your+Image&font=Montserrat"
          alt="Placeholder Image"
          class="object-cover w-full h-full"
        ></img>
      </div> */}
      {/* <!-- Right: Login Form --> */}
      <div class="lg:p-36 md:p-52 sm:20 p-8 w-full lg:w-1/2">
        <h1 class="mt-10 text-center text-2xl font-bold leading-9 tracking-tight text-gray-900">
          Login
        </h1>
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
          <div class="mb-6 text-green-400">
            <a href="#" class="hover:underline">
              Forgot Password?
            </a>
          </div>
          {/* <!-- Login Button --> */}
          <button
            type="submit"
            class="flex w-full justify-center rounded-md bg-green-600 px-3 py-1.5 text-sm font-semibold leading-6 text-white shadow-sm hover:bg-green-600 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-green-600"
          >
            Login
          </button>
        </form>
        {/* <!-- Sign up  Link --> */}
        <div class="mt-6 text-green-400 text-center">
          <a href="#" class="hover:underline">
            Sign up Here
          </a>
        </div>
      </div>
    </div>
  );
};

export default Login;
