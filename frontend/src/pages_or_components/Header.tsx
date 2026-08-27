import '../App.css'

function Header() {
    const goToDocs = () => {
        window.location.href= '/docs'
    }

    const goToGen = () => {
        window.location.href= '/'
    }
    return (
        <div className='flex justify-center items-center bg-blue-800 h-20 
        border-b border-blue-900'>
            <div className="flex space-x-8 py-4">
                <button className="text-lg px-6 py-3 text-white 
                transition-all duration-300 hover:text-2xl hover:bg-blue-700 
                rounded-lg hover:cursor-pointer" onClick={goToGen}>
                    Home
                </button>
                <button className="text-lg px-6 py-3 text-white 
                transition-all duration-300 hover:text-2xl hover:bg-blue-700 
                rounded-lg hover:cursor-pointer" onClick={goToDocs}>
                    Docs
                </button>
            </div>
        </div>

    )
}

export default Header