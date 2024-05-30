import React from "react";
// import { useState } from "react";

const Contact = (props) => {
  return (
    <section class="bg-white">
      <div class="py-8 lg:py-16 px-4 mx-auto max-w-screen-md">
        <h2 class="mb-4 text-4xl tracking-tight text-center">Contact Us</h2>
        <p class="mb-8 lg:mb-16 font-light text-center sm:text-xl">
          Got a technical issue? Want to send feedback about a beta feature?
          Need details about our Business plan? Let us know.
        </p>
        <div class="lg:p-28 md:p-36 sm:20 p-4 w-full shadow-md border border-gray-200 rounded-lg">
          <form action="#" class="space-y-8">
            <div>
              <label for="email" class="block mb-2 text-sm font-medium">
                Your email
              </label>
              <input
                type="email"
                id="email"
                class="w-full border border-gray-300 rounded-md py-2 px-3 focus:outline-none focus:border-green-400"
                placeholder="name@flowbite.com"
                required
              ></input>
            </div>
            <div>
              <label for="subject" class="block mb-2 text-sm font-medium ">
                Subject
              </label>
              <input
                type="text"
                id="subject"
                class="w-full border border-gray-300 rounded-md py-2 px-3 focus:outline-none focus:border-green-400"
                placeholder="Let us know how we can help you"
                required
              ></input>
            </div>
            <div class="sm:col-span-2">
              <label for="message" class="block mb-2 text-sm font-medium ">
                Your message
              </label>
              <textarea
                id="message"
                rows="6"
                class="w-full border border-gray-300 rounded-md py-2 px-3 focus:outline-none focus:border-green-400"
                placeholder="Leave a comment..."
              ></textarea>
            </div>
            <button
              type="submit"
              class="py-3 px-5 text-sm font-medium text-center text-white rounded-lg bg-primary-green sm:w-fit hover:bg-primary-800 focus:ring-4 focus:outline-none focus:ring-primary-300 dark:bg-primary-600 dark:hover:bg-primary-700 dark:focus:ring-primary-800"
            >
              Send message
            </button>
          </form>
        </div>
      </div>
    </section>
  );
};

export default Contact;
