import { createContext, useEffect, useRef, useState } from "react"
import { useNavigate } from "react-router-dom";
import { motion, useScroll, useTransform } from "motion/react";


export const AboutRefContext = createContext<React.RefObject<HTMLElement | null>>({ current: null, })

export default function HomePage() {

    const navigate = useNavigate()
    const aboutRef = useRef<HTMLDivElement | null>(null)
    const { scrollYProgress : aboutScrollY } = useScroll({
        target: aboutRef,
        offset: ["start end", "end start"],

    })
    const aboutY = useTransform(aboutScrollY, [0.05, 0.175], [200, 0])
    const aboutOpacity = useTransform(aboutScrollY, [0.1, 0.275], [0, 1])
    const about2Y = useTransform(aboutScrollY, [0.06, 0.21], [200, 0])
    const about2Opacity = useTransform(aboutScrollY, [0.11, 0.25], [0, 1])

    const featuresRef = useRef<HTMLDivElement | null>(null)
    const { scrollYProgress : featuresScrollY } = useScroll({
        target: featuresRef,
        offset: ["start end", "end start"]
    })

    const contactRef = useRef<HTMLDivElement | null>(null)
    const { scrollYProgress : contactScrollY } = useScroll({
        target: contactRef,
        offset: ["start end", "end start"]
    })

    // useEffect(() => {
    //     if (location.hash) {
    //     const el = document.querySelector(location.hash);
    //         if (el) {
    //             el.scrollIntoView({ behavior: "smooth" });
    //         }
    //     }
    // }, [location]);

    const [ link, setLink ] = useState<string>("")
    const handleSubmit = (e: React.FormEvent) => {   
        e.preventDefault()

        var trimmedLink = link.trim()
        if (!trimmedLink || !trimmedLink.startsWith('https://open.spotify.com/')){
            trimmedLink = 'https://open.spotify.com/'
        } 

        try {
            new URL(trimmedLink);
        } catch {
            trimmedLink = 'https://open.spotify.com/';
        }
        
        // sanitize minimally here, do more on generator page
        navigate(`/recsGenerator?link=${encodeURIComponent(trimmedLink)}`)
    }

    return (
        <div>
            <section className="h-[80vh] w-full flex flex-col justify-center items-center">
                <div className="mb-4 text-5xl mx-[2%] text-center">
                    <h1>Find new songs that you <i>actually</i> like.</h1>
                </div>
                <form className='max-w-[1600px] mx-[10%] text-black flex flex-row' onSubmit={handleSubmit}>
                    <div className='flex justify-center items-center'>
                    <input 
                    className='flex px-2 mx-2 w-auto max-w-600 min-w-100 bg-white'
                    type='text'
                    value={link}
                    // we don't want an error --> will prevent redirect. so, only sanitize enough to redirect here
                    // TODO also, can just take onsubmit if we want, dont need to constantly update (minimal optimizations)
                    onChange={e => setLink(e.target.value)}
                    placeholder='Enter a spotify link here'
                    />
                    </div>
                    <div className='flex justify-center items-center'>
                    <button type='submit' className='px-2 mx-2 rounded-md w-30 border-black border-1 bg-gray-300 hover:bg-gray-200 active:bg-gray-100'>Generate</button>
                    </div>
                </form>
            </section>

            {/* about section */}
            <section ref={aboutRef}
                id="About" className="pt-5 pb-10 text-center flex items-center justify-center h-auto bg-[#1c263f]"    
            >
                <div className="flex flex-col items-center justify-center">
                    <motion.h1 className="text-[clamp(1.75rem,6vw,3rem)] mb-5 border-white border-b-2 pb-5 mx-[10%]"
                        style={{ y: aboutY, opacity: aboutOpacity }}
                    >About</motion.h1>
                    <motion.p className="text-[clamp(1.25rem,3vw,1.5rem)] mx-[10%]"
                        style={{ y: about2Y, opacity: about2Opacity }}
                    >Give us a spotify link, and we'll give you suggested artists, albums, or songs, backed by others and ML models. Made by music listeners, for music listeners.</motion.p>
                </div>
            </section>

            {/* features/APIs section */}
            {/* different tools usable -- lastfm (artist, song, album),  */}
            <section ref={featuresRef}
                id="Features" className="pt-5 pb-10 text-center flex items-center justify-center h-auto w-full"    
            >
                <div className="flex flex-col items-center justify-center w-full">
                    <motion.h1 className="text-[clamp(1.75rem,6vw,3rem)] mb-5 border-white border-b-2 pb-5 mx-[10%]"
                        style={{ y: aboutY, opacity: aboutOpacity }}
                    >Features</motion.h1>

                    {/* deezer <-- spotify web api --> (discogs/(music brainz/acoustic brainz)) --> lastfm */}
                    {/* TODO coolors hover for desc of banner */}
                    {/* TODO sizes properly + hover onclick for docs */}
                    <div className={`w-full flex items-center justify-center ${window.innerWidth < 1250 ? "flex-col" : "flex-row"}`}>
                        <div className={`text-center flex flex-col justify-center items-center rounded-lg border-white border-1 px-8 mx-5 ${window.innerWidth < 1250 ? "max-w-200 min-w-100 mx-10 min-h-40 text-xl my-4" : "min-w-50 min-h-40 py-5 text-xl"}`}><h1 className="text-3xl">Spotify Web API</h1><p>Returns artist name, album name(optional), and song name(optional).</p></div>
                        <div className={`text-center flex flex-col justify-center items-center rounded-lg border-white border-1 px-8 mx-5 ${window.innerWidth < 1250 ? "max-w-200 min-w-100 mx-10 min-h-40 text-xl my-4" : "min-w-50 min-h-40 py-5 text-xl"}`}><h1 className="text-3xl">Last.fm API</h1><p>Recommends songs/artists/albums based on <i>tags</i>.</p></div>
                        <div className={`text-center flex flex-col justify-center items-center rounded-lg border-white border-1 px-8 mx-5 ${window.innerWidth < 1250 ? "max-w-200 min-w-100 mx-10 min-h-40 text-xl my-4" : "min-w-50 min-h-40 py-5 text-xl"}`}><h1 className="text-3xl">Discogs API</h1><p>Returns genres/styles to query Last.fm API with.</p></div>
                        <div className={`text-center flex flex-col justify-center items-center rounded-lg border-white border-1 px-8 mx-5 ${window.innerWidth < 1250 ? "max-w-200 min-w-100 mx-10 min-h-40 text-xl my-4" : "min-w-50 min-h-40 py-5 text-xl"}`}><h1 className="text-3xl">MusicBrainz/AcousticBrainz API</h1><p>Using ML models, returns genres/styles to query Last.fm API with.</p></div>
                        <div className={`text-center flex flex-col justify-center items-center rounded-lg border-white border-1 px-8 mx-5 ${window.innerWidth < 1250 ? "max-w-200 min-w-100 mx-10 min-h-40 text-xl my-4" : "min-w-50 min-h-40 py-5 text-xl"}`}><h1 className="text-3xl">Deezer API</h1><p>Recommends songs/artists/albums based on <i>similar artists</i>.</p></div>
                    </div>
                </div>
            </section>

            {/* Contact section */}
            <section ref={contactRef}
                id="How" className="pt-5 pb-10 text-center flex items-center justify-center h-auto "    
            >
                <div className="flex flex-col items-center justify-center">
                    <motion.h1 className="text-[clamp(1.75rem,6vw,3rem)] mb-5 border-white border-b-2 pb-5 mx-[10%]"
                        style={{ y: aboutY, opacity: aboutOpacity }}
                    >Contact</motion.h1>
                    <motion.p className="text-[clamp(1.25rem,3vw,1.5rem)] mx-[10%]"
                        style={{ y: about2Y, opacity: about2Opacity }}
                    >Contact me at <a href="mailto:liamtamh@gmail.com" target="_blank" className="text-blue-500 hover:text-blue-400 active:text-blue-300 underline">my email</a>, <a href="https://www.linkedin.com/in/liam-huynh-91aa1a1a1/" target="_blank" className="text-blue-500 hover:text-blue-400 active:text-blue-300 underline">my linkedin</a>, or <a href="https://github.com/huynhli" target="_blank" className="text-blue-500 hover:text-blue-400 active:text-blue-300 underline">my github</a>.</motion.p>
                </div>
            </section>
        </div>
    )
}