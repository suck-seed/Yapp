'use client';

import React from 'react';
import { useState } from 'react';
import { motion } from 'framer-motion';
import Image from 'next/image';
import logo from '../../assets/images/yapLogo.png';
import { FaEye, FaEyeSlash } from "react-icons/fa";
import { FcGoogle } from "react-icons/fc";
import Link from 'next/link';

export default function SignIn(){

    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [showPassword, setShowPassword] = useState(false);
    const [formError, setFormError] = useState('');
    const [focusedField, setFocusedField] = useState('');

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setFormError('');
        
        try {
            const res = await fetch('api/auth/login',{
                method: "POST",
                headers: { "content-type": "application/json"},
                body: JSON.stringify({ email, password }),
            });

            const data = await res.json();

            if(res.ok){
                alert('Login successful');
                localStorage.setItem('token', data.token);
                window.location.href = "/home";
            } else {
                setFormError(data.message);
            }
        } catch(error){
            console.error(error); // Log error for debugging
            setFormError('An error occurred. Please try again later.');
        }
    };

    const getInputBorder = (value: string, isFocused: boolean) => {
        if (value || isFocused) return 'border-[#0077d4]';
        return 'border-[#dcd9d3]';
    };

    return(
        <div className="min-h-screen flex flex-col bg-center font-[SF_Pro_Display] bg-transparent bg-[radial-gradient(#000000_1px,#e5e5f7_1px)] bg-[length:30px_30px] ">

            <main className="flex-1 flex justify-center items-center m-8 rounded-[20px]">
                <motion.div className="flex flex-col justify-center bg-white rounded-3xl p-8 z-10 border-3 z-1"
                initial={{opacity:0, y:40}}
                animate={{opacity:1, y:0}}
                transition={{duration:0.6, ease: 'easeOut'}}>

                    <div className='flex justify-start items-start'>
                        <a href="/get-started" className={`flex flex-row justify-start items-start relative -inset-2`}><Image className='w-7' src={logo} alt='Yapp Logo'/></a>
                    </div>

                    <section className="flex flex-col justify-center m-8 mt-2">

                        <h1 className="text-xl font-semibold font-[SF_Pro_Rounded] mt-4 text-[#1e1e1e]">Yapp — Connect. Collaborate. Communicate.</h1>

                        <p className='text-xl font-medium text-[#B6B09F] mb-8 tracking-wide'>Login to your Yapp account</p>

                        <form onSubmit={handleSubmit} className="flex flex-col justify-center flex-1 text-[#1e1e1e] font-[SF_Pro_Rounded]">

                            <div className='flex flex-col mb-4 gap-5'>
                                <button
                                  type="submit"
                                  className="bg-white text-black py-3 rounded-lg text-base w-full cursor-pointer rounded-3xl border-3 border-[#dcd9d3]"
                                >
                                    <div className='flex flex-row justify-start gap-2 ml-4'>
                                        <FcGoogle className='flex justify-center items-center' size={28}/>
                                        <p className='justify-center items-center text-lg ml-19'>Continue with Google</p>
                                    </div>

                                </button>

                                <div className="flex items-center my-2 w-full">
                                    <div className="flex-grow h-px bg-gray-400 opacity-35" />
                                    <div className="flex-grow h-px bg-gray-400 opacity-35" />
                                </div>
                            </div>

                            <div className="flex flex-col gap-1 mb-8">
                                <label className="text-sm text-[#73726e] font-medium">
                                  Email {!email && <span className="text-red-600">*</span>}
                                </label>
                                <input
                                  type="email"
                                  className={`rounded-lg px-2 py-3 pl-3 w-full bg-white text-black font-light border-3  font-[SF_Pro_Display] 
                                    ${getInputBorder(email, focusedField === 'email')} focus:outline-none focus:border-[#0077d4]`}
                                  value={email}
                                  onChange={(e) => setEmail(e.target.value)}
                                  onFocus={() => setFocusedField('email')}
                                  onBlur={() => setFocusedField('')}
                                  required
                                />
                            </div>

                            <div className="flex flex-col gap-1 mb-4">

                                <div className='flex flex-row items-center'>
                                    <label className="text-sm font-medium text-[#73726e] flex-1">
                                      Password {!password && <span className="text-red-600">*</span>}
                                    </label>
                                    <label className='text-xs text-blue'>
                                        <a href="/forget-password" className="text-[#0077d4]">Forget Password?</a>
                                    </label>
                                </div>

                                <div className="relative">
                                    <input
                                      type={showPassword ? 'text' : 'password'}
                                      className={`rounded-lg px-2 py-3 pl-3 w-full pr-10 bg-white text-black font-light border-3 font-[SF_Pro_Display] ${getInputBorder(password, focusedField === 'password')} focus:outline-none focus:border-[#0077d4]`}
                                      value={password}
                                      onChange={(e) => setPassword(e.target.value)}
                                      onFocus={() => setFocusedField('password')}
                                      onBlur={() => setFocusedField('')}
                                      required
                                    />
                                    <button
                                      type="button"
                                      onClick={() => setShowPassword(!showPassword)}
                                      className="absolute right-4 top-4 cursor-pointer"
                                    >
                                        {showPassword ? (
                                            <FaEye size={24} />
                                        ) : (
                                            <FaEyeSlash size={24} />
                                        )}
                                    </button>
                                </div>
                            </div>


                            <div className="flex items-center gap-2 mb-4 ">
                                <input
                                  type="checkbox"
                                  className="w-4 h-4 accent-[#0077d4]"
                                />
                                <label className="text-sm font-medium">Remember Me</label>
                            </div>

                            {formError && <p className="text-red-600 text-sm mb-4">{formError}</p>}

                            <div className='flex flex-col gap-2'>
                                <button
                                  type="submit"
                                  className="bg-[#2383E2] text-xl text-white py-3 rounded-lg text-base w-full cursor-pointer hover:bg-[#0077d4] font-medium"
                                >
                                  Login
                                </button>                                
                            </div>

                            <p className="flex justify-center text-sm mt-4 text-[#1e1e1e]">
                                Don&apos;t have an account?{' '}
                                <Link href="/signup" className="text-[#0077d4] ml-2">Sign Up</Link>
                            </p>
                        </form>
                    </section>
                </motion.div>
            
        </main>

        {/* <div className='w-full h-full '>
              <Image className='w-auto h-auto' src={bg} alt='communication'/>
        </div> */}
    </div>
    );
}