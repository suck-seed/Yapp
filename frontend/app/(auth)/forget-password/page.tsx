'use client'

import { useState } from 'react';
import { motion } from 'framer-motion';
import { MdLockReset } from "react-icons/md";
import { IoIosArrowBack } from "react-icons/io";
import Link from 'next/link';
// import Image from 'next/image';
// import logo from '@/app/assets/images/yapLogo.png'
// import gradientBG from '../../assets/images/gradient.png';

export default function ForgetPassword(){
    const [email, setEmail] = useState('');
    const [formError, setFormError] = useState('');
    const [focusedField, setFocusedField] = useState('');

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
    
        setFormError('');
        alert('Form submitted successfully!');
    };

    const getInputBorder = (value: string, isFocused: boolean) => {
        if (value || isFocused) return 'border-[#0077d4]';
        return 'border-[#dcd9d3]';
    };


    return (
        <div className="min-h-screen flex flex-col bg-center font-[SF_Pro_Display] bg-transparent bg-[radial-gradient(#000000_1px,#e5e5f7_1px)] bg-[length:30px_30px] ">

            <main className='flex-1 flex justify-center items-center m-8 rounded-[20px]'>

                <motion.div className='flex flex-col justify-center bg-white rounded-3xl p-8 z-1 border-3'
                initial={{opacity:0, y:40}}
                animate={{opacity:1, y:0}}
                transition={{duration:0.6, ease: 'easeOut'}}>

                    {/* <div className='flex justify-start items-start'>
                        <a href="/get-started" className={`flex flex-row justify-start items-start relative -inset-2`}><Image className='w-5' src={logo} alt='Yapp Logo'/></a>
                    </div> */}

                    <div className='flex justify-center items-center'>
                        <MdLockReset size={80} color="#73726e" />
                    </div>

                    <section className='flex flex-col justify-center m-4 items-center'>

                        <h1 className='text-xl font-semibold font-[SF_Pro_Rounded] text-[#1e1e1e]'>Forgot your password?</h1>
                        <p className='text-xl font-medium text-[#a7a7a7] mb-8 font-[SF_Pro_Rounded]'>
                            Enter your email to reset your password.
                        </p>
                        <form onSubmit={handleSubmit} className='flex flex-col jsutify-center items-center flex-1 text-[#1e1e1e] font-[SF_Pro_Rounded]'>
                            <div className='flex flex-col gap-1 mb-6'>
                            <label className="text-sm text-[#73726e] font-medium">
                                  Email {!email && <span className="text-red-600">*</span>}
                                </label>
                                <input
                                  type="email"
                                  className={`rounded-lg px-2 py-3 w-100 bg-white text-black font-light border ${getInputBorder(email, focusedField === 'email')} focus:outline-none focus:border-[#0077d4]`}
                                  value={email}
                                  onChange={(e) => setEmail(e.target.value)}
                                  onFocus={() => setFocusedField('email')}
                                  onBlur={() => setFocusedField('')}
                                  required
                                />
                            </div>

                            <div className="flex flex-col mb-4">
                                <button 
                                    type='submit'
                                    className='bg-[#2383E2] text-white py-3 rounded-lg text-lg w-100 cursor-pointer hover:bg-[#0077d4] font-medium'
                                >
                                    Send Code
                                </button>
                            </div>
                            <div className="flex items-center my-2 w-full">
                                <div className="flex-grow h-px bg-gray-400 opacity-35" />
                                <div className="flex-grow h-px bg-gray-400 opacity-35" />
                            </div>
                        </form>

                        <Link href="/signin">
                            <div className='flex justify-center items-center gap-4 mt-2'>
                            <IoIosArrowBack size={24} color="#1e1e1e" />
                            <p className='text-[#1e1e1e]'>Back to Login</p>
                        </div>
                        </Link>
                        
                        

                        {formError && <p className="text-red-600 text-sm mb-4">{formError}</p>}
                    </section>
                </motion.div>
            </main>
        </div>
    )
}