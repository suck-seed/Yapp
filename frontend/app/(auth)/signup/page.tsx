'use client';

import { useState } from 'react';
import { motion } from 'framer-motion';
import Image from 'next/image';
import logo from '../../assets/images/yapLogo.png';
import Link from 'next/link';
import { FaEye, FaEyeSlash } from "react-icons/fa";


export default function SignUp() {
  const [email, setEmail] = useState('');
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [agreeToTerms, setAgreeToTerms] = useState(false);
  const [formError, setFormError] = useState('');
  const [focusedField, setFocusedField] = useState('');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (!agreeToTerms) {
      setFormError('You must agree to the Terms, Privacy Policy, and Fees.');
      return;
    }

    if (password !== confirmPassword) {
      setFormError('Passwords do not match.');
      return;
    }

    setFormError('');
    alert('Form submitted successfully!');
  };

  const getInputBorder = (value: string, isFocused: boolean) => {
    if (value || isFocused) return 'border-[#0077d4]';
      return 'border-[#dcd9d3]';
  };

  return (
    <div className="min-h-screen flex flex-col bg-center font-[SF_Pro_Display] bg-transparent bg-[radial-gradient(#000000_1px,#e5e5f7_1px)] bg-[length:30px_30px] ">

      <main className="flex-1 flex justify-center items-center m-8 rounded-[20px]">

        <motion.div className="flex flex-col justify-center bg-white w-[545px] rounded-3xl p-8 z-10 border-3 z-1"
          initial={{opacity:0, y:40}}
          animate={{opacity:1, y:0}}
          transition={{duration:0.6, ease: 'easeOut'}}
        >
          <div className='flex justify-start items-start'>
            <a href="/get-started" className={`flex flex-row justify-start items-start relative -inset-2`}><Image className='w-7' src={logo} alt='Yapp Logo'/></a>
          </div> 

          <section className="flex flex-col justify-center m-8 mt-2 mb-2">
            <h1 className="text-xl font-semibold font-[SF_Pro_Rounded] mt-4 text-[#1e1e1e] w-full">Yapp — Connect. Collaborate. Communicate.</h1>

            <p className='text-xl font-medium text-[#B6B09F] mb-8 tracking-wide'>Register into Yapp account</p>

            <form onSubmit={handleSubmit} className="flex flex-col justify-center flex-1 text-[#1e1e1e] font-[SF_Pro_Rounded]">

              {/** Email Field */}
              <div className="flex flex-col gap-1 mb-8 ">
                <label className="text-sm text-[#73726e] font-medium">
                  Email {!email && <span 
                  className="text-red-600">*</span>}
                </label>
                <input
                  type="email"
                  className={`rounded-lg px-2 py-3 pl-3 w-full bg-white text-black font-light border-3 font-[SF_Pro_Display] 
                    ${getInputBorder(email, focusedField === 'email')} focus:outline-none focus:border-[#0077d4]`}
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  onFocus={() => setFocusedField('email')}
                  onBlur={() => setFocusedField('')}
                  required
                />
              </div>

              {/** Username Field */}
              <div className="flex flex-col gap-1 mb-8 ">
                <label className="text-sm text-[#73726e] font-medium">
                  Username {!username && <span className="text-red-600">*</span>}
                </label>
                <input
                  type="text"
                  className={`rounded-lg px-2 py-3 pl-3 w-full bg-white text-black font-light border-3 font-[SF_Pro_Display] ${getInputBorder(username, focusedField === 'username')} focus:outline-none focus:border-[#0077d4]`}
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  onFocus={() => setFocusedField('username')}
                  onBlur={() => setFocusedField('')}
                  required
                />
              </div>

              {/** Password Field */}
              <div className="flex flex-col gap-1 mb-8">
                <label className="text-sm text-[#73726e] font-medium">
                  Password {!password && <span className="text-red-600">*</span>}
                </label>
                <div className="relative">
                  <input
                    type={showPassword ? 'text' : 'password'}
                    className={`rounded-lg px-2 py-3 pl-3 w-full pr-10 bg-white text-black font-light border-3 ${getInputBorder(password, focusedField === 'password')} focus:outline-none focus:border-[#0077d4]`}
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    onFocus={() => setFocusedField('password')}
                    onBlur={() => setFocusedField('')}
                    required
                  />
                  <button
                    type="button"
                    onClick={() => setShowPassword(!showPassword)}
                    className="absolute right-3 top-3.5 cursor-pointer"
                  >
                    {showPassword ? (
                      <FaEye size={24} />
                      ) : (
                      <FaEyeSlash size={24} />
                    )}
                  </button>
                </div>
              </div>

              {/** Confirm Password Field */}
              <div className="flex flex-col gap-1 mb-4">
                <label className="text-sm text-[#73726e] font-medium">
                  Re-type Password {!confirmPassword && <span className="text-red-600">*</span>}
                </label>
                <div className="relative">
                <input
                    type={showConfirmPassword ? 'text' : 'password'}
                    className={`rounded-lg px-2 py-3 pr-10 w-full bg-white text-black font-light border-3 ${getInputBorder(confirmPassword, focusedField === 'confirm')} focus:outline-none focus:border-[#0077d4]`}
                    value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                    onFocus={() => setFocusedField('confirm')}
                    onBlur={() => setFocusedField('')}
                    required
                  />
                  <button
                    type="button"
                    onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                    className="absolute right-3 top-3.5 cursor-pointer"
                  >
                    {showPassword ? (
                      <FaEye size={24} />
                      ) : (
                      <FaEyeSlash size={24} />
                    )}
                  </button>
                </div>
              </div>

              <div className="flex items-center gap-2 mb-4">
                <input
                  type="checkbox"
                  checked={agreeToTerms}
                  onChange={(e) => setAgreeToTerms(e.target.checked)}
                  className="w-4 h-4 accent-[#0077d4]"
                />
                <label className="text-sm font-medium">I agree to all Term, Privacy Policy and Fees</label>
              </div>

              {formError && <p className="text-red-600 text-sm mb-4">{formError}</p>}

              <button
                type="submit"
                className="bg-[#2383E2] text-white py-3 rounded-lg text-lg w-full cursor-pointer hover:bg-[#0077d4] font-medium"
              >
                Sign Up
              </button>

              <p className="text-sm mt-2 text-[#1e1e1e] flex justify-center items center gap-2">
              Already have an account? 
              <Link href="/signin" className="text-[#1371FF]">Log In</Link>
              </p>
            </form>

            
          </section>
        </motion.div>
      </main>
    </div>
  );
}
