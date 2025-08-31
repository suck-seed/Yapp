'use client';
import React, { useState, useEffect, useCallback } from 'react';
import { motion } from 'framer-motion';
import Image from 'next/image';
import logo from '../../assets/images/yap final logo.png';
import { FaEye, FaEyeSlash, FaSpinner } from 'react-icons/fa';
import { FcGoogle } from 'react-icons/fc';
import Link from 'next/link';
import { useSignIn, useUser } from '@clerk/nextjs';
import { useRouter } from 'next/navigation';

// Type definitions for Clerk errors
interface ClerkError {
  code: string;
  message: string;
  longMessage?: string;
  meta?: Record<string, unknown>;
}

interface ClerkAPIError {
  errors: ClerkError[];
  clerkTraceId?: string;
}

// Type guard to check if error is a Clerk API error
const isClerkAPIError = (error: unknown): error is ClerkAPIError => {
  return (
    typeof error === 'object' &&
    error !== null &&
    'errors' in error &&
    Array.isArray((error as ClerkAPIError).errors) &&
    (error as ClerkAPIError).errors.length > 0
  );
};

// Form validation utilities
const validateEmail = (email: string): boolean => {
  const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
  return emailRegex.test(email.trim());
};

const validatePassword = (password: string): boolean => {
  return password.length >= 8;
};

// Main component
export default function SignIn() {
  const { signIn, setActive, isLoaded } = useSignIn();
  const { isSignedIn } = useUser();
  const router = useRouter();

  // Form state
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [rememberMe, setRememberMe] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  const [formError, setFormError] = useState('');
  const [focusedField, setFocusedField] = useState('');
  const [fieldErrors, setFieldErrors] = useState<{ email?: string; password?: string }>({});
  const [announceError, setAnnounceError] = useState('');
  const [isLoading, setIsLoading] = useState(false);


  useEffect(() => {
    if (isLoaded && isSignedIn) {
      router.replace('/home'); // Already signed-in users go straight to home
    }
  }, [isLoaded, isSignedIn, router]);

  // Load remember me preference on component mount
  useEffect(() => {
    const rememberMeCookie = document.cookie
      .split('; ')
      .find(row => row.startsWith('remember_me='));
    if (rememberMeCookie) {
      setRememberMe(true);
    }
  }, []);

  const getErrorMessage = useCallback((error: unknown): string => {
    if (isClerkAPIError(error)) {
      const errorCode = error.errors[0]?.code;
      const errorMessage = error.errors[0]?.message;

      switch (errorCode) {
        case 'form_identifier_not_found':
          return 'No account found with this email address';
        case 'form_password_incorrect':
          return 'Incorrect password. Please try again.';
        case 'too_many_requests':
          return 'Too many failed attempts. Please try again later.';
        case 'session_exists':
          return 'You are already signed in';
        case 'identifier_already_signed_in':
          return 'This account is already signed in from another session';
        case 'form_password_pwned':
          return 'This password has been found in a data breach. Please use a different password.';
        case 'form_param_nil':
          return 'Please fill in all required fields';
        case 'form_password_length_too_short':
          return 'Password must be at least 8 characters long';
        case 'form_identifier_invalid':
          return 'Please enter a valid email address';
        case 'strategy_for_user_invalid':
          return 'Invalid verification method. If you used Google for Sign Up use it for Sign In too.';
        default:
          // return errorCode + ': ' + errorMessage || 'An unexpected error occurred. Please try again.';
          return errorMessage || 'Sign In failed. Please try again.';
      }
    }

    if (error instanceof Error) {
      return error.message || 'An unexpected error occurred';
    }

    if (typeof error === 'string') {
      return error;
    }

    return 'An unexpected error occurred. Please try again.';
  }, []);

  const signInWithCredentials = async (email: string, password: string, rememberMe: boolean) => {
    if (!isLoaded) return;

    setIsLoading(true);
    try {
      const result = await signIn.create({
        identifier: email.trim().toLowerCase(),
        password,
      });

      if (result.status === 'complete') {
        await setActive({
          session: result.createdSessionId,
          beforeEmit: () => {
            if (rememberMe) {
              document.cookie = `remember_me=true; max-age=${30 * 24 * 60 * 60}; path=/; secure; samesite=strict`;
            }
          }
        });

        router.push('/home');
      } else if (
        result.status === 'needs_identifier' ||
        result.status === 'needs_first_factor' ||
        result.status === 'needs_second_factor' ||
        result.status === 'needs_new_password'
      ) {
        setFormError('Additional verification steps required');
      } else {
        setFormError('Sign In failed. Please try again.');
      }
    } catch (error: unknown) {
      const errorMessage = getErrorMessage(error);
      setFormError(errorMessage);
      setAnnounceError(errorMessage);
    } finally {
      setIsLoading(false);
    }
  };

  const signInWithGoogle = async () => {
    if (!isLoaded) return;

    setIsLoading(true);
    try {
      await signIn.authenticateWithRedirect({
        strategy: 'oauth_google',
        redirectUrl: `/signin/sso-callback`,
        redirectUrlComplete: `/home`,
      });
    } catch (error: unknown) {
      const errorMessage = getErrorMessage(error) || 'Google sign-in failed';
      setFormError(errorMessage);
      setAnnounceError(errorMessage);
    } finally {
      setIsLoading(false);
    }
  };

  // Real-time validation
  const validateField = useCallback((field: string, value: string) => {
    const errors = { ...fieldErrors };

    if (field === 'email') {
      if (value && !validateEmail(value)) {
        errors.email = 'Please enter a valid email address';
      } else {
        delete errors.email;
      }
    }

    if (field === 'password') {
      if (value && !validatePassword(value)) {
        errors.password = 'Password must be at least 8 characters long';
      } else {
        delete errors.password;
      }
    }

    setFieldErrors(errors);
  }, [fieldErrors]);

  const handleEmailChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setEmail(value);
    validateField('email', value);
    if (formError) setFormError('');
  };

  const handlePasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setPassword(value);
    validateField('password', value);
    if (formError) setFormError('');
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (isLoading) return;

    // Validate all fields
    const errors: { email?: string; password?: string } = {};

    if (!email.trim()) {
      errors.email = 'Email is required';
    } else if (!validateEmail(email)) {
      errors.email = 'Please enter a valid email address';
    }

    if (!password) {
      errors.password = 'Password is required';
    } else if (!validatePassword(password)) {
      errors.password = 'Password must be at least 8 characters long';
    }

    if (Object.keys(errors).length > 0) {
      setFieldErrors(errors);
      const errorMessage = 'Please fix the form errors before submitting';
      setFormError(errorMessage);
      setAnnounceError(errorMessage);
      return;
    }

    setFormError('');
    setFieldErrors({});

    await signInWithCredentials(email, password, rememberMe);
  };

  const handleGoogleSignIn = async () => {
    if (isLoading) return;

    setFormError('');
    await signInWithGoogle();
  };

  const getInputBorder = (value: string, isFocused: boolean, hasError: boolean) => {
    if (hasError) return 'border-red-500';
    return value || isFocused ? 'border-[#0077d4]' : 'border-[#dcd9d3]';
  };

  const isFormValid = email.trim() && password && validateEmail(email) && validatePassword(password) && Object.keys(fieldErrors).length === 0;

  return (
    <div className="min-h-screen flex flex-col bg-center font-[SF_Pro_Display] bg-transparent bg-[radial-gradient(#000000_1px,#e5e5f7_1px)] bg-[length:30px_30px]">
      {/* Screen reader announcements */}
      <div
        role="status"
        aria-live="polite"
        aria-atomic="true"
        className="sr-only"
      >
        {announceError}
      </div>

      <main className="flex-1 flex justify-center items-center m-8 rounded-[20px]">
        <motion.div
          className="flex flex-col justify-center bg-white rounded-3xl p-8 z-10 border-3 max-w-md w-full"
          initial={{ opacity: 0, y: 40 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6, ease: 'easeOut' }}
        >
          <div className="flex justify-start items-start">
            <Link
              href="/get-started"
              className="flex flex-row justify-start items-start relative -inset-2"
              aria-label="Go back to get started page"
            >
              <Image className="w-7" src={logo} alt="Yapp Logo" priority />
            </Link>
          </div>

          <section className="flex flex-col justify-center m-8 mt-2">
            <h1 className="text-xl font-semibold font-[SF_Pro_Rounded] mt-4 text-[#1e1e1e]">
              Yapp — Connect. Collaborate. Communicate.
            </h1>
            <p className="text-xl font-medium text-[#B6B09F] mb-8 tracking-wide">
              Sign In to your Yapp account
            </p>

            <form
              onSubmit={handleSubmit}
              className="flex flex-col justify-center text-[#1e1e1e] font-[SF_Pro_Rounded]"
              noValidate
            >
              <div className="flex flex-col mb-4 gap-5">
                <button
                  type="button"
                  onClick={handleGoogleSignIn}
                  disabled={isLoading}
                  className="bg-white text-black py-3 text-base w-full cursor-pointer rounded-3xl border-3 border-[#dcd9d3] hover:border-[#0077d4] transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                  aria-label="Sign in with Google"
                >
                  <div className="flex flex-row justify-start gap-2 ml-4">
                    <span className="flex justify-center items-center">
                      {isLoading ? <FaSpinner className="animate-spin" size={20} /> : <FcGoogle size={28} />}
                    </span>
                    <p className="justify-center items-center text-lg ml-19">
                      {isLoading ? 'Signing in...' : 'Continue with Google'}
                    </p>
                  </div>
                </button>

                <div className="flex items-center my-2 w-full" role="separator" aria-label="or">
                  <div className="flex-grow h-px bg-gray-400 opacity-35" />
                  <span className="px-2 text-gray-500 text-sm">or</span>
                  <div className="flex-grow h-px bg-gray-400 opacity-35" />
                </div>
              </div>

              <div className="flex flex-col gap-1 mb-4">
                <label
                  htmlFor="email"
                  className="text-sm text-[#73726e] font-medium"
                >
                  Email {!email && <span className="text-red-600" aria-label="required">*</span>}
                </label>
                <input
                  id="email"
                  name="email"
                  type="email"
                  autoComplete="email"
                  className={`rounded-lg px-2 py-3 pl-3 w-full bg-white text-black font-light border-3 font-[SF_Pro_Display] ${getInputBorder(email, focusedField === 'email', !!fieldErrors.email)} focus:outline-none focus:border-[#0077d4] transition-colors`}
                  value={email}
                  onChange={handleEmailChange}
                  onFocus={() => setFocusedField('email')}
                  onBlur={() => setFocusedField('')}
                  aria-describedby={fieldErrors.email ? 'email-error' : undefined}
                  aria-invalid={!!fieldErrors.email}
                  disabled={isLoading}
                  required
                />
                {fieldErrors.email && (
                  <span id="email-error" className="text-red-600 text-sm mt-1" role="alert">
                    {fieldErrors.email}
                  </span>
                )}
              </div>

              <div className="flex flex-col gap-1 mb-4">
                <div className="flex flex-row items-center">
                  <label
                    htmlFor="password"
                    className="text-sm font-medium text-[#73726e] flex-1"
                  >
                    Password {!password && <span className="text-red-600" aria-label="required">*</span>}
                  </label>
                  <Link
                    href="/forgot-password"
                    className="text-xs text-[#0077d4] hover:underline"
                    tabIndex={isLoading ? -1 : 0}
                  >
                    Forgot Password?
                  </Link>
                </div>
                <div className="relative">
                  <input
                    id="password"
                    name="password"
                    type={showPassword ? 'text' : 'password'}
                    autoComplete="current-password"
                    className={`rounded-lg px-2 py-3 pl-3 w-full pr-12 bg-white text-black font-light border-3 font-[SF_Pro_Display] ${getInputBorder(password, focusedField === 'password', !!fieldErrors.password)} focus:outline-none focus:border-[#0077d4] transition-colors`}
                    value={password}
                    onChange={handlePasswordChange}
                    onFocus={() => setFocusedField('password')}
                    onBlur={() => setFocusedField('')}
                    aria-describedby={fieldErrors.password ? 'password-error' : undefined}
                    aria-invalid={!!fieldErrors.password}
                    disabled={isLoading}
                    required
                  />
                  <button
                    type="button"
                    onClick={() => setShowPassword(!showPassword)}
                    className="absolute right-4 top-4 cursor-pointer hover:text-[#0077d4] transition-colors disabled:cursor-not-allowed"
                    aria-label={showPassword ? 'Hide password' : 'Show password'}
                    disabled={isLoading}
                    tabIndex={isLoading ? -1 : 0}
                  >
                    {showPassword ? <FaEye size={20} /> : <FaEyeSlash size={20} />}
                  </button>
                </div>
                {fieldErrors.password && (
                  <span id="password-error" className="text-red-600 text-sm mt-1" role="alert">
                    {fieldErrors.password}
                  </span>
                )}
              </div>

              <div className="flex items-center gap-2 mb-4">
                <input
                  id="remember-me"
                  type="checkbox"
                  className="w-4 h-4 accent-[#0077d4] focus:ring-2 focus:ring-[#0077d4] focus:ring-offset-1"
                  checked={rememberMe}
                  onChange={(e) => setRememberMe(e.target.checked)}
                  disabled={isLoading}
                />
                <label htmlFor="remember-me" className="text-sm font-medium cursor-pointer">
                  Remember Me for 30 days
                </label>
              </div>

              {formError && (
                <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded-lg">
                  <p className="text-red-600 text-sm" role="alert" aria-live="polite">
                    {formError}
                  </p>
                </div>
              )}

              {/* CAPTCHA Widget - Clerk will inject the CAPTCHA */}
              <div id="clerk-captcha" className="mb-4"></div>

              <div className="flex flex-col gap-2">
                <button
                  type="submit"
                  disabled={isLoading || !isFormValid}
                  className="bg-[#2383E2] text-white py-3 rounded-lg text-base w-full cursor-pointer hover:bg-[#0077d4] font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed focus:ring-2 focus:ring-[#0077d4] focus:ring-offset-2"
                  aria-describedby="Sign-In-button-description"
                >
                  {isLoading ? (
                    <div className="flex items-center justify-center gap-2">
                      <FaSpinner className="animate-spin" size={16} />
                      Signing in...
                    </div>
                  ) : (
                    'Sign In'
                  )}
                </button>
                <div id="Sign In-button-description" className="sr-only">
                  {!isFormValid && 'Please fill out all required fields correctly to enable Sign In'}
                </div>
              </div>

              <p className="flex justify-center text-sm mt-4 text-[#1e1e1e]">
                Don&apos;t have an account?{' '}
                <Link
                  href="/signup"
                  className="text-[#0077d4] ml-2 hover:underline focus:underline"
                  tabIndex={isLoading ? -1 : 0}
                >
                  Sign Up with email
                </Link>
              </p>
            </form>
          </section>
        </motion.div>
      </main>
    </div>
  );
}