import "../styles/Header.css"
function Header() {

    return (
      <>
        <header>
            <img src="/dictionary.png" alt="logo" />
            <h1>Write it down</h1>
            <button className="themeBtn">
                <img src="/day.png" alt="theme" />
                {/* <img src="/night.png" alt="theme" /> */}
            </button>
        </header>
      </>
    )
}

export default Header;