import React from "react";

export default function Navigation() {
    return (
        <header className="bg-green-600">
            <nav className="mx-auto flex max-w-full items-center justify-between p-6">
               <div className="flex flex-1 justify-start gap-4">
                    <a href="#" className="-m-1.5 p-1.5">
                        logo
                    </a>
                </div> 
                <div className="flex flex-shrink-0 flex-grow justify-start items-center gap-4 px-2">
                    <a href="#" className="text-lg font-semibold leading-6 text-white">Technologien</a>
                    <a href="#" className="text-lg font-semibold leading-6 text-white">Aanvragen</a>
                    <a href="#" className="text-lg font-semibold leading-6 text-white">Prijzen</a>
                    <a href="#" className="text-lg font-semibold leading-6 text-white">Handleiding</a>
                </div>
                <div className="flex flex-shrink-0 flex-grow justify-end items-center gap-4">
                    <a href="#" className="text-lg font-semibold leading-6 text-white">login</a>
                    <button className="inline-block rounded-md border border-transparent bg-green-500 px-8 py-2 text-center font-medium text-white hover:bg-green-400">Account aanmaken</button>
                    <button className="inline-block rounded-md border border-transparent bg-green-500 px-8 py-2 text-center font-medium text-white hover:bg-green-400">jouw omgeving</button>
                </div>
            </nav>   
        </header>
        
    )
}