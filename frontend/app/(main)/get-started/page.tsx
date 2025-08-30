'use client';

// import { Space_Grotesk } from 'next/font/google';
import React from 'react';
import Image from 'next/image';
import { useRef, useEffect} from 'react';
import gsap from 'gsap';
import right from '../../assets/images/Right.png';
import yapLogo from '../../assets/images/yapLogo.png';
import Link from 'next/link';
import comm from '../../assets/images/communication.png';
import { motion, Variants } from 'framer-motion';

const containerVariants = {
  hidden: {},
  show: {
    transition: {
      staggerChildren: 0,
    },
  },
};

const lineVariants:Variants = {
  hidden: { y: 60, opacity: 0 },
  show: {
    y: 0,
    opacity: 1,
    transition: {
      ease: 'easeIn',
      duration: 1,
    },
  },
};

export default function CoverPage() {
  const paraRef = useRef(null);

  useEffect(() => {
    if (paraRef.current) {
      gsap.fromTo(paraRef.current, 
        { y: 50, opacity: 0 }, 
        { y:0, opacity: 1, duration: 1, ease: 'power3.out' }
      );
    }
  }, []);

  return (
    
    <div className="bg-[#d6c4a4] flex flex-col min-h-screen bg-[#f4f4f4]">

      <header className="flex justify-left items-center">
        <div className='ml-8 mt-8 flex justify-center items-center'>
            <Image className="w-7" src={yapLogo} alt="Yapp logo" />
        </div>
      </header>

      <main className="flex-1 flex flex-col justify-between items-center m-8 rounded-[20px] gap-20">

        <div className="flex flex-1 flex-col justify-start items-center h-full w-full p-10">

          <motion.div
            className="text-7xl text-center mt-[16px] mb-[32px] font-bold leading-tight tracking-tight text-[#1e1e1e] font-[Space_Grotesk]"
            variants={containerVariants}
            initial="hidden"
            animate="show"
          >
            <motion.div variants={lineVariants}>CONNECT.</motion.div>
            <motion.div variants={lineVariants}>COLLABORATE.</motion.div>
            <motion.div variants={lineVariants}>COMMUNICATE.</motion.div>
          </motion.div>

            {/* <motion.h1
              className="text-7xl text-center mt-[4rem] mb-[4rem] font-bold leading-tight tracking-tight text-[#1e1e1e] font-[Heuvel_Grotesk_Demo]"
              initial={{ y: 60, opacity: 0 }}
              animate={{ y: 0, opacity: 1 }}
              transition={{ duration: 1, ease: 'easeOut' }}
            >
              CONNECT.<br />
              COLLABORATE.<br />
              COMMUNICATE.
            </motion.h1> */}

            <div className="flex flex-col justify-center items-center z-10">
                <motion.p 
                  className="text-3xl text-[#1e1e1e] font-[SF_Pro_Rounded] mb-[32px] font-base tracking-tight"
                  initial={{opacity:0, y:50}}
                  animate={{opacity:1, y:0}}
                  transition={{duration:1, ease:'easeOut'}}
                >
                The fast, friendly way to stay connected with the people who matter most.
                </motion.p>

                <Link href="/signin">
                    <motion.button
                      className="flex items-center gap-2 bg-[#1e1e1e] text-white px-8 py-4 rounded-lg hover:bg-black cursor-pointer"
                      initial={{ opacity: 0, y: 50 }}
                      animate={{ opacity: 1, y: 0 }}
                      transition={{ duration: 1, ease: "easeOut" }}
                      >
                      <motion.div
                        initial={{ opacity: 0 }}
                        animate={{ opacity: 1 }}
                        transition={{ delay: 0.2, duration: 1 }}
                        className="flex items-center gap-2"
                      >
                        <p className="text-2xl font-[SF_Pro_Rounded]">Get Started</p>
                        <Image className="w-5 h-5" src={right} alt="right-arrow" priority />
                      </motion.div>
                    </motion.button>

                </Link>
            </div>
            <div className='absolute top-[42rem] z-0'>
              <Image className='w-auto h-auto' src={comm} alt='communication'/>
            </div>
        </div>
      </main>
    </div>
  );
};


