import '../App.css'
import { AnimatePresence, easeInOut, motion } from "motion/react"
import { useState } from 'react'
import MotionHashLink from './MotionHashLink'

function Header() {
    // window location href is for new tab (re-renders otherwise) -> also for scripts. if not script, use <a target="_blank"></a>
    // useNavigate is for scripts
    // use Link or NavLink --> NavLink gives class="active"

    const [navBarOpen, setNavBarOpen] = useState<Boolean>(false)

    const similarVariant = {
        initial: {width: 100, opacity: 1},
        animate: {width: 0, opacity: 0},
        hover: {width: 100, opacity: 1}
    }

    const logoVariant = {
        initial: { x: 120 },
        animate: { x: 0, transition: {duration: 0.5, ease: easeInOut, delay: 0.4}},
        hover: { x: 120, transition: {duration: 0.3, ease: easeInOut}}
    }
    const ongsVariant = {
        initial: { x: 0, opacity: 1 },
        animate: { x: -100, opacity: 0, transition: {duration: 0.5, ease: easeInOut, delay: 0.4}},
        hover: { x: 0, opacity: 1, transition: {duration: 0.3, ease: easeInOut}}
    }

    return (
        <header 
            className='
                flex justify-between items-center relative 
                h-20 w-auto mx-[4%] mt-[10px] pb-[10px]
                border-b-1 border-white'
        >
            <MotionHashLink
                className='relative flex flex-row h-full text-4xl text-white font-electrolize' 
                to='/'
                smooth
                whileHover="hover" 
                initial="initial"
                animate="animate"
                key="logo-link"
                >
                 {/* text slide left */}
                <motion.div className='absolute top-[16px] right-20 h-full origin-right'
                    transition={{width: {duration: 0.5, ease: easeInOut, delay: 0.2}, opacity: {duration: 0.3, delay: 0.2}}}
                    variants={similarVariant}
                >
                    Similar
                </motion.div>
                
                {/* img + txt right */}
                <motion.div 
                    // TODO bg same colour as header
                    className='flex justify-center items-center h-full z-0 bg-[#141b2a]'
                    variants={logoVariant}
                >
                    <img src='/similarSongsLogo.png' className='h-full'/>
                    <motion.div className='w-30 justify-left z-[-10]'
                        variants={ongsVariant}
                    >ongs</motion.div>
                </motion.div>
            </MotionHashLink>
            
            <button className='flex justify-center items-center flex-col w-20 h-20 group relative'
                onClick={() => setNavBarOpen(prev => !prev)}
            >
                {/* hamburg */}
                <div className={`bg-white w-[50%] h-1 m-1 transition-transform duration-600 ${navBarOpen ? "-rotate-270 translate-y-[11px] translate-x-[9px]" : ""}`}/>
                <div className={`bg-white w-[50%] h-1 m-1 transition-transform duration-600 ${navBarOpen ? "-rotate-270 -translate-x-[9px]" : ""}`}/>
                <div className={`bg-white w-[50%] h-1 m-1 transition-transform duration-600 ${navBarOpen ? "-rotate-270 -translate-y-[13px]" : ""}`}/>
            </button>
            
            {/* navbar */}
            { navBarOpen && 
                <AnimatePresence>
                    <motion.nav className='absolute w-full bottom-[-80.2vh] left-0 h-[80vh] px-[10%] bg-[#141b2a] flex flex-col text-4xl z-0'
                        initial={{ y: "-100%", opacity: 0}}
                        animate={{ y: 0, opacity: 1 }}
                        transition={{y: {duration: 0.2, ease: "easeInOut"}, opacity: {duration: 0.3, ease: "easeIn"}}}
                        exit={{}}
                    >
                        <MotionHashLink className='flex-1 px-14 mb-2 mt-20 flex items-center hover:font-serif hover:italic text-zinc-300 hover:text-white'
                            onClick={() => (setNavBarOpen(prev => !prev))}
                            to='/recsGenerator'
                            smooth
                            initial={{ y: -100, opacity: 0}}
                            animate={{ y: 0, opacity: 1}}
                            transition={{duration: 0.3, ease: easeInOut, delay: 0.3 }}
                        >
                            <h1>Recommend a song!</h1>
                        </MotionHashLink>

                        <MotionHashLink className='z-10 flex-1 px-14 my-2 flex items-center hover:font-serif hover:italic text-zinc-300 hover:text-white'
                            onClick={() => (setNavBarOpen(prev => !prev))}
                            to='/#About'
                            smooth
                            initial={{ y: -100, opacity: 0}}
                            animate={{ y: 0, opacity: 1}}
                            transition={{duration: 0.3, ease: easeInOut, delay: 0.35 }}
                        >
                            About
                        </MotionHashLink>

                        <MotionHashLink className='z-10 flex-1 px-14 my-2 flex items-center hover:font-serif hover:italic text-zinc-300 hover:text-white'
                            onClick={() => (setNavBarOpen(prev => !prev))}
                            to='/#Features'
                            smooth
                            initial={{ y: -100, opacity: 0}}
                            animate={{ y: 0, opacity: 1}}
                            transition={{duration: 0.3, ease: easeInOut, delay: 0.4 }}
                        >
                            Features
                        </MotionHashLink>

                        <MotionHashLink className='z-10 flex-1 px-14 my-2 flex items-center hover:font-serif hover:italic text-zinc-300 hover:text-white'
                            onClick={() => (setNavBarOpen(prev => !prev))}
                            to='/#Contacts'
                            smooth
                            initial={{ y: -100, opacity: 0}}
                            animate={{ y: 0, opacity: 1}}
                            transition={{duration: 0.3, ease: easeInOut, delay: 0.45 }}
                        >
                            Contacts
                        </MotionHashLink>
                    </motion.nav>
                </AnimatePresence>
            }
            
        </header>

    )
}

export default Header