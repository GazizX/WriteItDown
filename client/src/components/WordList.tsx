import WordItem from "./WordItem";
import { useQuery } from "@tanstack/react-query";
import { BASE_URL } from "../App";
export type Word = {
    id: number;
    body: string;
    translation: string;
}
function WordList() {
  const { data: words} = useQuery<Word[]>({
		queryKey: ["words"],
		queryFn: async () => {
			try {
				const res = await fetch(BASE_URL + "/words");
				const data = await res.json();

				if (!res.ok) {
					throw new Error(data.error || "Something went wrong");
				}
				return data || [];
			} catch (error) {
				console.log(error);
			}
		},
	});
    return (
      <>
        {words?.map((word) => (
					<WordItem key={word.id} word={word} />
				)).reverse()}
      </>
    )
}

export default WordList;