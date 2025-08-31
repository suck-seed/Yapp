'use client';
import React, { useState, useCallback, useEffect } from 'react';
import { motion } from 'framer-motion';
import Image from 'next/image';
import logo from '../../assets/images/yap final logo.png';
import Link from 'next/link';
import { FaEye, FaEyeSlash, FaSpinner, FaCheck, FaTimes } from "react-icons/fa";
import { useSignUp, useUser } from '@clerk/nextjs';
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

// Password strength enum
enum PasswordStrength {
  WEAK = 'weak',
  FAIR = 'fair',
  GOOD = 'good',
  STRONG = 'strong'
}

// Form validation utilities
const validateEmail = (email: string): boolean => {
  const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
  return emailRegex.test(email.trim());
};

const validateUsername = (username: string): boolean => {
  return /^[a-zA-Z0-9_]{3,20}$/.test(username.trim());
};

const checkPasswordStrength = (password: string): PasswordStrength => {
  let score = 0;

  if (password.length >= 8) score++;
  if (/[a-z]/.test(password)) score++;
  if (/[A-Z]/.test(password)) score++;
  if (/[0-9]/.test(password)) score++;
  if (/[^a-zA-Z0-9]/.test(password)) score++;
  if (password.length >= 12) score++;

  if (score <= 2) return PasswordStrength.WEAK;
  if (score <= 3) return PasswordStrength.FAIR;
  if (score <= 4) return PasswordStrength.GOOD;
  return PasswordStrength.STRONG;
};

const validatePassword = (password: string): { isValid: boolean; errors: string[] } => {
  const errors: string[] = [];

  if (password.length < 8) {
    errors.push('At least 8 characters long');
  }
  if (!/[a-z]/.test(password)) {
    errors.push('One lowercase letter');
  }
  if (!/[A-Z]/.test(password)) {
    errors.push('One uppercase letter');
  }
  if (!/[0-9]/.test(password)) {
    errors.push('One number');
  }
  if (!/[^a-zA-Z0-9]/.test(password)) {
    errors.push('One special character');
  }

  return { isValid: errors.length === 0, errors };
};

// Main component
export default function SignUp() {
  const { signUp, setActive, isLoaded } = useSignUp();
  const { isSignedIn } = useUser();
  const router = useRouter();

  useEffect(() => {
    if (isLoaded && isSignedIn) {
      router.replace('/home'); // Already signed-in users go straight to home
    }
  }, [isLoaded, isSignedIn, router]);

  // Form state
  const [email, setEmail] = useState('');
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [agreeToTerms, setAgreeToTerms] = useState(false);
  const [formError, setFormError] = useState('');
  const [focusedField, setFocusedField] = useState('');
  const [fieldErrors, setFieldErrors] = useState<{ [key: string]: string }>({});
  const [announceError, setAnnounceError] = useState('');
  const [showVerification, setShowVerification] = useState(false);
  const [verificationCode, setVerificationCode] = useState('');
  const [isLoading, setIsLoading] = useState(false);



  const getErrorMessage = useCallback((error: unknown): string => {
    if (isClerkAPIError(error)) {
      const errorCode = error.errors[0]?.code;
      const errorMessage = error.errors[0]?.message;

      switch (errorCode) {
        case 'form_identifier_exists':
          return 'An account with this email already exists. Please sign in instead.';
        case 'form_password_pwned':
          return 'This password has been found in a data breach. Please use a different password.';
        case 'form_password_length_too_short':
          return 'Password must be at least 8 characters long';
        case 'form_password_no_uppercase_letter':
          return 'Password must contain at least one uppercase letter';
        case 'form_password_no_lowercase_letter':
          return 'Password must contain at least one lowercase letter';
        case 'form_password_no_number':
          return 'Password must contain at least one number';
        case 'form_password_no_special_char':
          return 'Password must contain at least one special character';
        case 'form_identifier_invalid':
          return 'Please enter a valid email address';
        case 'form_param_nil':
          return 'Please fill in all required fields';
        case 'too_many_requests':
          return 'Too many attempts. Please try again later.';
        case 'form_username_invalid':
          return 'Username can only contain letters, numbers, and underscores';
        case 'form_username_exists':
          return 'This username is already taken. Please choose another.';
        default:
          return errorMessage || 'Registration failed. Please try again.';
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

  //Temporarily Removed the need for username in sign up flow

  const signUpWithCredentials = async (email: string, password: string) => {
    if (!isLoaded) return;

    setIsLoading(true);
    try {
      const result = await signUp.create({
        emailAddress: email.trim().toLowerCase(),
        password,
        // username: username.trim()
      });

      if (result.status === 'complete') {
        await setActive({ session: result.createdSessionId });
        router.push('/home');
      } else if (result.status === 'missing_requirements') {
        await result.prepareEmailAddressVerification({ strategy: 'email_code' });
        setShowVerification(true);
      } else {
        setFormError('Additional verification steps required');
      }
    } catch (error: unknown) {
      const errorMessage = getErrorMessage(error);
      setFormError(errorMessage);
      setAnnounceError(errorMessage);
    } finally {
      setIsLoading(false);
    }
  };

  const verifyEmail = async (code: string) => {
    if (!isLoaded || !signUp) return;

    setIsLoading(true);
    try {
      const result = await signUp.attemptEmailAddressVerification({ code });

      if (result.status === 'complete') {
        await setActive({ session: result.createdSessionId });
        router.push('/home');
      } else {
        setFormError('Verification incomplete. Please try again.');
      }
    } catch (error: unknown) {
      const errorMessage = getErrorMessage(error);
      setFormError(errorMessage);
      setAnnounceError(errorMessage);
    } finally {
      setIsLoading(false);
    }
  };

  // Real-time validation
  const validateField = useCallback((field: string, value: string) => {
    const errors = { ...fieldErrors };

    switch (field) {
      case 'email':
        if (value && !validateEmail(value)) {
          errors.email = 'Please enter a valid email address';
        } else {
          delete errors.email;
        }
        break;
      case 'username':
        if (value && !validateUsername(value)) {
          errors.username = 'Username must be 3-20 characters, letters, numbers, and underscores only';
        } else {
          delete errors.username;
        }
        break;
      case 'password':
        const passwordValidation = validatePassword(value);
        if (value && !passwordValidation.isValid) {
          errors.password = 'Password requirements not met';
        } else {
          delete errors.password;
        }
        break;
      case 'confirmPassword':
        if (value && value !== password) {
          errors.confirmPassword = 'Passwords do not match';
        } else {
          delete errors.confirmPassword;
        }
        break;
    }

    setFieldErrors(errors);
  }, [fieldErrors, password]);

  const handleInputChange = (field: string, value: string) => {
    switch (field) {
      case 'email':
        setEmail(value);
        break;
      case 'username':
        setUsername(value);
        break;
      case 'password':
        setPassword(value);
        if (confirmPassword) {
          validateField('confirmPassword', confirmPassword);
        }
        break;
      case 'confirmPassword':
        setConfirmPassword(value);
        break;
    }

    validateField(field, value);
    if (formError) setFormError('');
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (isLoading) return;

    // Validate all fields
    const errors: { [key: string]: string } = {};

    if (!email.trim()) {
      errors.email = 'Email is required';
    } else if (!validateEmail(email)) {
      errors.email = 'Please enter a valid email address';
    }

    if (!username.trim()) {
      errors.username = 'Username is required';
    } else if (!validateUsername(username)) {
      errors.username = 'Username must be 3-20 characters, letters, numbers, and underscores only';
    }

    if (!password) {
      errors.password = 'Password is required';
    } else {
      const passwordValidation = validatePassword(password);
      if (!passwordValidation.isValid) {
        errors.password = 'Password requirements not met';
      }
    }

    if (!confirmPassword) {
      errors.confirmPassword = 'Please confirm your password';
    } else if (confirmPassword !== password) {
      errors.confirmPassword = 'Passwords do not match';
    }

    if (!agreeToTerms) {
      errors.terms = 'You must agree to all Terms, Privacy Policy and Fees';
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

    await signUpWithCredentials(email, password);
  };

  const handleVerificationSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!verificationCode.trim()) {
      setFormError('Please enter the verification code');
      return;
    }

    await verifyEmail(verificationCode);
  };

  const getInputBorder = (value: string, isFocused: boolean, hasError: boolean = false) => {
    if (hasError) return 'border-red-500';
    if (value || isFocused) return 'border-[#0077d4]';
    return 'border-[#dcd9d3]';
  };

  const passwordStrength = password ? checkPasswordStrength(password) : null;
  const passwordValidation = password ? validatePassword(password) : { isValid: false, errors: [] };

  const getPasswordStrengthColor = (strength: PasswordStrength) => {
    switch (strength) {
      case PasswordStrength.WEAK: return 'bg-red-500';
      case PasswordStrength.FAIR: return 'bg-yellow-500';
      case PasswordStrength.GOOD: return 'bg-blue-500';
      case PasswordStrength.STRONG: return 'bg-green-500';
    }
  };

  const getPasswordStrengthWidth = (strength: PasswordStrength) => {
    switch (strength) {
      case PasswordStrength.WEAK: return 'w-1/4';
      case PasswordStrength.FAIR: return 'w-2/4';
      case PasswordStrength.GOOD: return 'w-3/4';
      case PasswordStrength.STRONG: return 'w-full';
    }
  };

  const isFormValid = email.trim() && username.trim() && password &&
    confirmPassword && password === confirmPassword && agreeToTerms &&
    validateEmail(email) && validateUsername(username) &&
    passwordValidation.isValid && Object.keys(fieldErrors).length === 0;

  if (showVerification) {
    return (
      <div className="min-h-screen flex flex-col bg-center font-[SF_Pro_Display] bg-transparent bg-[radial-gradient(#000000_1px,#e5e5f7_1px)] bg-[length:30px_30px]">
        <main className="flex-1 flex justify-center items-center m-8 rounded-[20px]">
          <motion.div
            className="flex flex-col justify-center bg-white w-[545px] rounded-3xl p-8 z-10 border-3"
            initial={{ opacity: 0, y: 40 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6, ease: 'easeOut' }}
          >
            <div className="text-center mb-8">
              <h1 className="text-xl font-semibold font-[SF_Pro_Rounded] text-[#1e1e1e] mb-4">
                Verify Your Email
              </h1>
              <p className="text-gray-600 mb-4">
                We&apos;ve sent a verification code to <strong>{email}</strong>
              </p>
              <p className="text-sm text-gray-500">
                Please check your email and enter the code below
              </p>
            </div>

            <form onSubmit={handleVerificationSubmit} className="space-y-4">
              <div>
                <label htmlFor="verification-code" className="block text-sm font-medium text-[#73726e] mb-2">
                  Verification Code
                </label>
                <input
                  id="verification-code"
                  type="text"
                  className="w-full px-4 py-3 border-3 border-[#dcd9d3] rounded-lg focus:border-[#0077d4] focus:outline-none text-center text-lg font-mono tracking-wider text-gray-900 placeholder-gray-400"
                  value={verificationCode}
                  onChange={(e) => setVerificationCode(e.target.value.replace(/\D/g, '').slice(0, 6))}
                  placeholder="000000"
                  maxLength={6}
                  disabled={isLoading}
                />
              </div>

              {formError && (
                <div className="p-3 bg-red-50 border border-red-200 rounded-lg">
                  <p className="text-red-600 text-sm" role="alert">
                    {formError}
                  </p>
                </div>
              )}

              <button
                type="submit"
                disabled={isLoading || verificationCode.length !== 6}
                className="w-full bg-[#2383E2] text-white py-3 rounded-lg font-medium hover:bg-[#0077d4] transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {isLoading ? (
                  <div className="flex items-center justify-center gap-2">
                    <FaSpinner className="animate-spin" size={16} />
                    Verifying...
                  </div>
                ) : (
                  'Verify Email'
                )}
              </button>
            </form>
          </motion.div>
        </main>
      </div>
    );
  }

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
          className="flex flex-col justify-center bg-white w-[545px] rounded-3xl p-8 z-10 border-3"
          initial={{ opacity: 0, y: 40 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6, ease: 'easeOut' }}
        >
          <div className='flex justify-start items-start'>
            <Link
              href="/get-started"
              className="flex flex-row justify-start items-start relative -inset-2"
              aria-label="Go back to get started page"
            >
              <Image className='w-7' src={logo} alt='Yapp Logo' priority />
            </Link>
          </div>

          <section className="flex flex-col justify-center m-8 mt-2 mb-2">
            <h1 className="text-xl font-semibold font-[SF_Pro_Rounded] mt-4 text-[#1e1e1e] w-full">
              Yapp — Connect. Collaborate. Communicate.
            </h1>
            <p className='text-xl font-medium text-[#B6B09F] mb-8 tracking-wide'>
              Register into Yapp account
            </p>

            <form
              onSubmit={handleSubmit}
              className="flex flex-col justify-center flex-1 text-[#1e1e1e] font-[SF_Pro_Rounded]"
              noValidate
            >
              {/** Email Field */}
              <div className="flex flex-col gap-1 mb-8">
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
                  className={`rounded-lg px-2 py-3 pl-3 w-full bg-white text-black font-light border-3 font-[SF_Pro_Display] 
                    ${getInputBorder(email, focusedField === 'email', !!fieldErrors.email)} focus:outline-none focus:border-[#0077d4] transition-colors`}
                  value={email}
                  onChange={(e) => handleInputChange('email', e.target.value)}
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

              {/** Username Field */}
              <div className="flex flex-col gap-1 mb-8">
                <label
                  htmlFor="username"
                  className="text-sm text-[#73726e] font-medium"
                >
                  Username {!username && <span className="text-red-600" aria-label="required">*</span>}
                </label>
                <input
                  id="username"
                  name="username"
                  type="text"
                  autoComplete="username"
                  className={`rounded-lg px-2 py-3 pl-3 w-full bg-white text-black font-light border-3 font-[SF_Pro_Display] ${getInputBorder(username, focusedField === 'username', !!fieldErrors.username)} focus:outline-none focus:border-[#0077d4] transition-colors`}
                  value={username}
                  onChange={(e) => handleInputChange('username', e.target.value)}
                  onFocus={() => setFocusedField('username')}
                  onBlur={() => setFocusedField('')}
                  aria-describedby={fieldErrors.username ? 'username-error' : 'username-help'}
                  aria-invalid={!!fieldErrors.username}
                  disabled={isLoading}
                  required
                />
                {fieldErrors.username ? (
                  <span id="username-error" className="text-red-600 text-sm mt-1" role="alert">
                    {fieldErrors.username}
                  </span>
                ) : (
                  <span id="username-help" className="text-gray-500 text-xs mt-1">
                    3-20 characters, letters, numbers, and underscores only
                  </span>
                )}
              </div>

              {/** Password Field */}
              <div className="flex flex-col gap-1 mb-8">
                <label
                  htmlFor="password"
                  className="text-sm text-[#73726e] font-medium"
                >
                  Password {!password && <span className="text-red-600" aria-label="required">*</span>}
                </label>
                <div className="relative">
                  <input
                    id="password"
                    name="password"
                    type={showPassword ? 'text' : 'password'}
                    autoComplete="new-password"
                    className={`rounded-lg px-2 py-3 pl-3 w-full pr-12 bg-white text-black font-light border-3 font-[SF_Pro_Display] ${getInputBorder(password, focusedField === 'password', !!fieldErrors.password)} focus:outline-none focus:border-[#0077d4] transition-colors`}
                    value={password}
                    onChange={(e) => handleInputChange('password', e.target.value)}
                    onFocus={() => setFocusedField('password')}
                    onBlur={() => setFocusedField('')}
                    aria-describedby="password-requirements"
                    aria-invalid={!!fieldErrors.password}
                    disabled={isLoading}
                    required
                  />
                  <button
                    type="button"
                    onClick={() => setShowPassword(!showPassword)}
                    className="absolute right-3 top-3.5 cursor-pointer hover:text-[#0077d4] transition-colors disabled:cursor-not-allowed"
                    aria-label={showPassword ? 'Hide password' : 'Show password'}
                    disabled={isLoading}
                    tabIndex={isLoading ? -1 : 0}
                  >
                    {showPassword ? <FaEye size={24} /> : <FaEyeSlash size={24} />}
                  </button>
                </div>

                {/* Password strength indicator */}
                {password && passwordStrength && (
                  <div className="mt-2 space-y-1">
                    <div className="flex items-center gap-2">
                      <div className="flex-1 h-2 bg-gray-200 rounded-full overflow-hidden">
                        <div className={`h-full transition-all duration-300 ${getPasswordStrengthColor(passwordStrength)} ${getPasswordStrengthWidth(passwordStrength)}`} />
                      </div>
                      <span className="text-xs font-medium capitalize text-gray-600">
                        {passwordStrength}
                      </span>
                    </div>
                  </div>
                )}

                {/* Password requirements */}
                <div id="password-requirements" className="mt-2 space-y-1">
                  {password && passwordValidation.errors.map((error, index) => (
                    <div key={index} className="flex items-center gap-1 text-xs">
                      <FaTimes className="text-red-500" size={10} />
                      <span className="text-red-600">{error}</span>
                    </div>
                  ))}
                  {password && passwordValidation.isValid && (
                    <div className="flex items-center gap-1 text-xs">
                      <FaCheck className="text-green-500" size={10} />
                      <span className="text-green-600">Password meets all requirements</span>
                    </div>
                  )}
                  {!password && (
                    <div className="text-xs text-gray-500 space-y-1">
                      <div>Password must contain:</div>
                      <div className="ml-2 space-y-0.5">
                        <div>• At least 8 characters long</div>
                        <div>• One lowercase letter</div>
                        <div>• One uppercase letter</div>
                        <div>• One number</div>
                        <div>• One special character</div>
                      </div>
                    </div>
                  )}
                </div>
              </div>

              {/** Confirm Password Field */}
              <div className="flex flex-col gap-1 mb-4">
                <label
                  htmlFor="confirmPassword"
                  className="text-sm text-[#73726e] font-medium"
                >
                  Re-type Password {!confirmPassword && <span className="text-red-600" aria-label="required">*</span>}
                </label>
                <div className="relative">
                  <input
                    id="confirmPassword"
                    name="confirmPassword"
                    type={showConfirmPassword ? 'text' : 'password'}
                    autoComplete="new-password"
                    className={`rounded-lg px-2 py-3 pl-3 w-full pr-12 bg-white text-black font-light border-3 font-[SF_Pro_Display] ${getInputBorder(confirmPassword, focusedField === 'confirm', !!fieldErrors.confirmPassword)} focus:outline-none focus:border-[#0077d4] transition-colors`}
                    value={confirmPassword}
                    onChange={(e) => handleInputChange('confirmPassword', e.target.value)}
                    onFocus={() => setFocusedField('confirm')}
                    onBlur={() => setFocusedField('')}
                    aria-invalid={!!fieldErrors.confirmPassword}
                    disabled={isLoading}
                    required
                  />
                  <button
                    type="button"
                    onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                    className="absolute right-3 top-3.5 cursor-pointer hover:text-[#0077d4] transition-colors disabled:cursor-not-allowed"
                    aria-label={showConfirmPassword ? 'Hide confirm password' : 'Show confirm password'}
                    disabled={isLoading}
                    tabIndex={isLoading ? -1 : 0}
                  >
                    {showConfirmPassword ? <FaEye size={24} /> : <FaEyeSlash size={24} />}
                  </button>
                </div>
                {fieldErrors.confirmPassword && (
                  <span id="confirmPassword-error" className="text-red-600 text-sm mt-1" role="alert">
                    {fieldErrors.confirmPassword}
                  </span>
                )}
                {confirmPassword && !fieldErrors.confirmPassword && confirmPassword === password && (
                  <div className="flex items-center gap-1 text-xs mt-1">
                    <FaCheck className="text-green-500" size={10} />
                    <span className="text-green-600">Passwords match</span>
                  </div>
                )}
              </div>

              <div className="flex items-center gap-2 mb-4">
                <input
                  id="terms"
                  type="checkbox"
                  checked={agreeToTerms}
                  onChange={(e) => setAgreeToTerms(e.target.checked)}
                  className="w-4 h-4 accent-[#0077d4] focus:ring-2 focus:ring-[#0077d4] focus:ring-offset-1"
                  aria-describedby={fieldErrors.terms ? 'terms-error' : undefined}
                  aria-invalid={!!fieldErrors.terms}
                  disabled={isLoading}
                  required
                />
                <label htmlFor="terms" className="text-sm font-medium cursor-pointer">
                  I agree to all{' '}
                  <Link href="/terms" className="text-[#0077d4] hover:underline" target="_blank">
                    Terms
                  </Link>
                  ,{' '}
                  <Link href="/privacy" className="text-[#0077d4] hover:underline" target="_blank">
                    Privacy Policy
                  </Link>
                  {' '}and{' '}
                  <Link href="/fees" className="text-[#0077d4] hover:underline" target="_blank">
                    Fees
                  </Link>
                </label>
              </div>
              {fieldErrors.terms && (
                <span id="terms-error" className="text-red-600 text-sm mb-4 block" role="alert">
                  {fieldErrors.terms}
                </span>
              )}

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
                  className="bg-[#2383E2] text-white py-3 rounded-lg text-lg w-full cursor-pointer hover:bg-[#0077d4] font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed focus:ring-2 focus:ring-[#0077d4] focus:ring-offset-2"
                  aria-describedby="signup-button-description"
                >
                  {isLoading ? (
                    <div className="flex items-center justify-center gap-2">
                      <FaSpinner className="animate-spin" size={16} />
                      Creating Account...
                    </div>
                  ) : (
                    'Sign Up'
                  )}
                </button>
                <div id="signup-button-description" className="sr-only">
                  {!isFormValid && 'Please fill out all required fields correctly to create account'}
                </div>
              </div>

              <p className="text-sm mt-2 text-[#1e1e1e] flex justify-center items-center gap-2">
                Already have an account?
                <Link
                  href="/signin"
                  className="text-[#1371FF] hover:underline focus:underline"
                  tabIndex={isLoading ? -1 : 0}
                >
                  Sign In
                </Link>
              </p>
            </form>
          </section>
        </motion.div>
      </main>
    </div>
  );
}