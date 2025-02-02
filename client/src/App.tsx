import "./App.css"
import Header from './components/Header'
import WordForm from './components/WordForm'
import WordList from './components/WordList'
export const BASE_URL = import.meta.env.MODE === "development" ? "http://localhost:5000/api" : "/api";
function App() {

  return (
    <>
      <Header/>
      <WordForm/>
      <WordList/>
    </>
  )
}

export default App
