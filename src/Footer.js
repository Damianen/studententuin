import React from "react";

export default function Footer() {
    return (
        <footer className="bg-green-600">
            <div className="mx-auto flex max-w-full items-center justify-between p-6">
                <div className="flex flex-1 justify-start items-center">
                    <a href="#" className="-m-1.5 p-1.5">
                        logo
                    </a>
                </div>
                <div className="flex flex-1 justify-end items-center gap-4 -px-2">
                    <a href="#" className="text-lg font-semibold leading-6 text-white">over ons</a>
                    <a href="#" className="text-lg font-semibold leading-6 text-white">contact</a>
                </div>
            </div>
        </footer>
    );
}