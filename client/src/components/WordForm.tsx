import { useState} from "react";
import { FiPlus } from "react-icons/fi";
import "../styles/WordForm.css"
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { BASE_URL } from "../App";
function WordForm() {
  const [word, setWord] = useState<string>("");

  const queryClient = useQueryClient();

  const { mutate: createWord} = useMutation({
		mutationKey: ["createWord"],
		mutationFn: async (e: React.FormEvent) => {
			e.preventDefault();
			try {
				const res = await fetch(BASE_URL + `/words`, {
					method: "POST",
					headers: {
						"Content-Type": "application/json",
					},
					body: JSON.stringify({ body: word }),
				});
				const data = await res.json();

				if (!res.ok) {
					throw new Error(data.error || "Something went wrong");
				}

				setWord("");
				return data;
			} catch (error: any) {
				throw new Error(error);
			}
		},
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["words"] });
		},
		onError: (error: any) => {
			alert(error.message);
		},
	});

  return (
    <form onSubmit={createWord}>
      <input
        type="text"
        placeholder="Введите слово"
        value={word}
        onChange={(e) => setWord(e.target.value)}
        className="wordInput"
      />
      <button type="submit" className="wordCreateBtn">
        <FiPlus size={20} />
      </button>
    </form>
  );
};

export default WordForm;