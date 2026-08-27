import '../App.css'

function Footer() {
    return (
        <div className='bg-blue-800 min-h-10 flex justify-center items-center'>
            <h1 className='text-white italic'>Thanks for stopping by! If you'd like to learn more about 
                this project or me, check out my github <a href='https://github.com/huynhli' target='_blank' className='underline text-blue'>
                here</a>.</h1>
        </div>
    )
}

export default Footer