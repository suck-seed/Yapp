'use client';

import { useState } from 'react';
import Image from 'next/image';
import logo from '../../assets/images/yap final logo.png';
import speechBubble from '../../assets/images/chat GPT image.png';
import see from '../../assets/images/see.png';
import hide from '../../assets/images/hide.png';
import google from '../../assets/images/google.png';
import gradientBG from '../../assets/images/gradient.png';

export default function SignIn(){

    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [showPassword, setShowPassword] = useState(false);
    const [formError, setFormError] = useState('');
    const [focusedField, setFocusedField] = useState('');

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
    
        setFormError('');
        alert('Form submitted successfully!');
    };


    const getInputBorder = (value: string, isFocused: boolean) => {
        if (value || isFocused) return 'border-black';
        return 'border-gray-400';
    };

    return(
        <div className="min-h-screen flex flex-col bg-[#d6c4a4] bg-cover bg-center min-h-screen" style={{backgroundImage:`url(${gradientBG.src})`}}>
            <header className="flex justify-left gap-2 items-center">
                <div className='ml-16 mt-8 gap-2 flex justify-center items-center'>
                    <Image className="w-8" src={logo} alt="Yapp logo" />
                    <h1 className='text-[#1e1e1e]  flex justify-left items-center text-3xl font-[Heuvel_Grotesk_Demo] font-semibold'>Yapp</h1>
                </div>
            </header>

            <main className="flex-1 flex justify-center items-center">
                <div className="flex gap-20 bg-white rounded-3xl p-2">
                    {/* <section className="flex flex-col items-center justify-center w-[580px] h-[630px] rounded-3xl bg-[#1f1f1f]">
                        <Image className="w-[480px]" src={speechBubble} alt="Speech bubble" />
                        <p className="text-white text-2xl font-[Heuvel_Grotesk_Demo] mt-4">Yapp — Connect. Collaborate. Communicate.</p>
                    </section> */}

                    <section className="flex flex-col justify-center m-16">
                        <h1 className="text-6xl font-bold mt-4 mb-4 font-[Heuvel_Grotesk_Demo] text-[#1e1e1e] tracking-wide">Login</h1>
                        <p className='text-xl font-medium text-[#1e1e1e] mb-16'>Hi, Welcome back!👋🏼</p>
                        <form onSubmit={handleSubmit} className="flex flex-col justify-center flex-1 text-[#1e1e1e] font-[SF_Pro_Rounded]">

                            <div className="flex flex-col gap-1 mb-8">
                                <label className="text-lg font-medium">
                                  Email {!email && <span className="text-red-600">*</span>}
                                </label>
                                <input
                                  type="email"
                                  className={`rounded-lg px-2 py-3 w-100 bg-white text-black font-light border ${getInputBorder(email, focusedField === 'email')} focus:outline-none focus:border-black`}
                                  value={email}
                                  onChange={(e) => setEmail(e.target.value)}
                                  onFocus={() => setFocusedField('email')}
                                  onBlur={() => setFocusedField('')}
                                  required
                                />
                            </div>

                            <div className="flex flex-col gap-1 mb-4">

                                <div className='flex flex-row items-center'>
                                    <label className="text-lg font-medium flex-1">
                                      Password {!password && <span className="text-red-600">*</span>}
                                    </label>
                                    <label className='text-sm text-blue'>
                                        <a href="/forget" className="text-[#1371FF]">Forget Password?</a>
                                    </label>
                                </div>

                                <div className="relative">
                                    <input
                                      type={showPassword ? 'text' : 'password'}
                                      className={`rounded-lg px-2 py-3 w-100 pr-10 bg-white text-black font-light border ${getInputBorder(password, focusedField === 'password')} focus:outline-none focus:border-black`}
                                      value={password}
                                      onChange={(e) => setPassword(e.target.value)}
                                      onFocus={() => setFocusedField('password')}
                                      onBlur={() => setFocusedField('')}
                                      required
                                    />
                                    <button
                                      type="button"
                                      onClick={() => setShowPassword(!showPassword)}
                                      className="absolute right-2 top-2.5 cursor-pointer"
                                    >
                                        <Image src={showPassword ? see : hide} alt="Toggle password" width={20} height={20} />
                                    </button>
                                </div>
                            </div>


                            <div className="flex items-center gap-2 mb-4">
                                <input
                                  type="checkbox"
                                  className="w-4 h-4"
                                />
                                <label className="text-sm font-medium">Remember Me</label>
                            </div>

                            {formError && <p className="text-red-600 text-sm mb-4">{formError}</p>}

                            <div className='flex flex-col gap-2'>
                                <button
                                  type="submit"
                                  className="bg-[#1f1f1f] text-xl text-white py-3 rounded-lg text-lg w-100 cursor-pointer hover:bg-black font-medium"
                                >
                                  Login
                                </button>

                                <div className="flex items-center my-2 w-100">
                                    <div className="flex-grow h-px bg-gray-400 opacity-50" />
                                    <span className="px-2 text-sm">OR</span>
                                    <div className="flex-grow h-px bg-gray-400 opacity-50" />
                                </div>


                                <button
                                  type="submit"
                                  className="bg-white text-black py-3 rounded-lg text-lg w-100 cursor-pointer rounded-3xl border"
                                >
                                    <div className='flex flex-row justify-center gap-2 items-center'>
                                        <Image className='w-6 h-6' src={google} alt='Google Logo'></Image>
                                        Sign in with Google
                                    </div>

                                </button>
                            </div>


                            <p className="flex justify-center text-sm mt-4 text-[#1e1e1e]">
                            Don&apos;thave an account?{' '}
                            <a href="/signup" className="text-[#1371FF] ml-2">Sign Up</a>
                            </p>
                        </form>
                    </section>
                </div>
            </main>
    </div>
    );
}