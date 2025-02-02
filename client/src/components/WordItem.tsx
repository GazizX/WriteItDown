import { FiTrash2 } from "react-icons/fi";
import { Word } from "./WordList";
import "../styles/WordItem.css"
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { BASE_URL } from "../App";
const WordItem = ({ word }: { word: Word }) => {
  const queryClient = useQueryClient();
  const { mutate: deleteWord} = useMutation({
		mutationKey: ["deleteWord"],
		mutationFn: async () => {
			try {
				const res = await fetch(BASE_URL + `/words/${word.id}`, {
					method: "DELETE",
				});
				const data = await res.json();
				if (!res.ok) {
					throw new Error(data.error || "Something went wrong");
				}
				return data;
			} catch (error) {
				console.log(error);
			}
		},

		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["words"] });
		},
	});


  return (
    <div className="word-item">
      <div className="word-body">
        <span className="word">{word.body}</span>
        <span className="translation">{word.translation}</span>
      </div>
      <button onClick={() => deleteWord()} className="delete-btn">
        <FiTrash2 size={18} />
      </button>
    </div>
  );
};

export default WordItem;
