import Image from 'next/image';
import chatBubble from '../../assets/images/chat bubble.png';
import right from '../../assets/images/Right.png';
import yapLogo from '../../assets/images/yap final logo.png';
import Link from 'next/link';
import gradientBG from '../../assets/images/gradient.png';


export default function CoverPage() {
  return (
    
    <div className="min-h-screen bg-[#d6c4a4] flex flex-col bg-cover bg-center min-h-screen" style={{backgroundImage:`url(${gradientBG.src})`}}>
      <header className="flex justify-left gap-2 items-center">
        <div className='ml-16 mt-8 gap-2 flex justify-center items-center'>
            <Image className="w-8" src={yapLogo} alt="Yapp logo" />
            <h1 className='text-[#1e1e1e]  flex justify-left items-center text-3xl font-[Heuvel_Grotesk_Demo] font-semibold'>Yapp</h1>
        </div>
      </header>

      <main className="flex-1 flex justify-between items-center ml-16 mt-8 rounded-[20px]">

        <div className="flex flex-1 flex-col justify-end h-[880px] pr-20">
            <h1 className="text-7xl mb-[14rem] font-bold leading-tight tracking-tight text-[#1e1e1e] font-[Heuvel_Grotesk_Demo]">
                CONNECT.<br />
                COLLABORATE.<br />
                COMMUNICATE.
            </h1>

            <div className="space-y-6">
                <p className="text-3xl text-[#1e1e1e] font-[SF_Pro_Rounded] font-base tracking-tight">
                The fast, friendly way to stay connected with the people who matter most.
                </p>

                <Link href="/signin">
                    <button className="flex items-center gap-2 bg-[#1e1e1e] text-white px-12 py-6 rounded-md hover:bg-black transition cursor-pointer">
                        <p className='text-2xl font-[SF_Pro_Rounded]'>
                            Get Started
                        </p>
                        <Image className="w-5 h-5" src={right} alt='right-arrow' />
                    </button>
                </Link>
            </div>
        </div>

        <div className="w-[692px] h-[880px] rounded-3xl flex flex-1 items-top justify-center bg-transparent">
          <Image
            src={chatBubble}
            alt="chat-bubble"
            className="w-180 h-220 object-contain"
          />
        </div>
      </main>
    </div>
  );
};
