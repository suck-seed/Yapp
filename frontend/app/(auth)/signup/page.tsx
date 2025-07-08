'use client';

import { useState } from 'react';
import Image from 'next/image';
import logo from '../../assets/images/yap final logo.png';
import speechBubble from '../../assets/images/chat GPT image.png';
import see from '../../assets/images/see.png';
import hide from '../../assets/images/hide.png';
import gradientBG from '../../assets/images/gradient.png';

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
    if (value || isFocused) return 'border-black';
    return 'border-gray-400';
  };

  

  return (
    <div className="min-h-screen flex flex-col bg-[#d6c4a4] bg-cover bg-center min-h-screen" style={{backgroundImage:`url(${gradientBG.src})`}}>
      <header className="flex justify-left gap-2 items-center">
        <div className='ml-16 mt-8 gap-2 flex justify-center items-center'>
            <Image className="w-8" src={logo} alt="Yapp logo" />
            <h1 className='text-[#1e1e1e]  flex justify-left items-center text-3xl font-[Heuvel_Grotesk_Demo] font-semibold'>Yapp</h1>
        </div>
      </header>

      <main className="flex-1 flex justify-center items-center">
        <div className="flex gap-20 bg-white rounded-3xl p-8">
          <section className="flex flex-col items-center justify-center w-[580px] h-[630px] rounded-3xl bg-[#1f1f1f]">
            <Image className="w-[480px]" src={speechBubble} alt="Speech bubble" />
            <p className="text-white text-2xl font-[Heuvel_Grotesk_Demo] mt-4">Yapp — Connect. Collaborate. Communicate.</p>
          </section>

          <section className="flex flex-col mr-10">
            <h1 className="text-3xl font-semibold mt-4 mb-4 font-[Heuvel_Grotesk_Demo] text-[#1e1e1e] ">Start Yapping Today</h1>
            <form onSubmit={handleSubmit} className="flex flex-col justify-center flex-1 text-[#1e1e1e] font-[SF_Pro_Rounded]">
              {/** Email Field */}
              <div className="flex flex-col gap-1 mb-8">
                <label className="text-lg font-medium">
                  Email {!email && <span className="text-red-600">*</span>}
                </label>
                <input
                  type="email"
                  className={`rounded-lg px-2 py-2 w-80 bg-white text-black font-light border ${getInputBorder(email, focusedField === 'email')} focus:outline-none focus:border-black`}
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  onFocus={() => setFocusedField('email')}
                  onBlur={() => setFocusedField('')}
                  required
                />
              </div>

              {/** Username Field */}
              <div className="flex flex-col gap-1 mb-8">
                <label className="text-lg font-medium">
                  Username {!username && <span className="text-red-600">*</span>}
                </label>
                <input
                  type="text"
                  className={`rounded-lg px-2 py-2 w-80 bg-white text-black font-light border ${getInputBorder(username, focusedField === 'username')} focus:outline-none focus:border-black`}
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  onFocus={() => setFocusedField('username')}
                  onBlur={() => setFocusedField('')}
                  required
                />
              </div>

              {/** Password Field */}
              <div className="flex flex-col gap-1 mb-8">
                <label className="text-lg font-medium">
                  Password {!password && <span className="text-red-600">*</span>}
                </label>
                <div className="relative">
                  <input
                    type={showPassword ? 'text' : 'password'}
                    className={`rounded-lg px-2 py-2 w-80 pr-10 bg-white text-black font-light border ${getInputBorder(password, focusedField === 'password')} focus:outline-none focus:border-black`}
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

              {/** Confirm Password Field */}
              <div className="flex flex-col gap-1 mb-4">
                <label className="text-lg font-medium">
                  Re-type Password {!confirmPassword && <span className="text-red-600">*</span>}
                </label>
                <div className="relative">
                <input
                    type={showConfirmPassword ? 'text' : 'password'}
                    className={`rounded-lg px-2 py-2 w-80 pr-10 bg-white text-black font-light border ${getInputBorder(confirmPassword, focusedField === 'confirm')} focus:outline-none focus:border-black`}
                    value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                    onFocus={() => setFocusedField('confirm')}
                    onBlur={() => setFocusedField('')}
                    required
                  />
                  <button
                    type="button"
                    onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                    className="absolute right-2 top-2.5 cursor-pointer"
                  >
                    <Image src={showConfirmPassword ? see : hide} alt="Toggle password" width={20} height={20} />
                  </button>
                </div>
              </div>

              <div className="flex items-center gap-2 mb-4">
                <input
                  type="checkbox"
                  checked={agreeToTerms}
                  onChange={(e) => setAgreeToTerms(e.target.checked)}
                  className="w-4 h-4"
                />
                <label className="text-sm font-light">I agree to all Term, Privacy Policy and Fees</label>
              </div>

              {formError && <p className="text-red-600 text-sm mb-4">{formError}</p>}

              <button
                type="submit"
                className="bg-[#1f1f1f] text-white py-3 rounded-lg text-lg w-80 cursor-pointer hover:bg-black"
              >
                Sign Up
              </button>

              <p className="text-sm mt-4 text-[#1e1e1e]">
              Already have an account?{' '}
              <a href="/signin" className="text-[#1371FF]">Log In</a>
            </p>
            </form>

            
          </section>
        </div>
      </main>
    </div>
  );
}
