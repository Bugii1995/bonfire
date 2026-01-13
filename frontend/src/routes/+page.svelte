<script>
	let sessionId = null;
	let question = null;
	let explanation = "";
	let status = "idle"; // idle | answering | finished
	let mastery = null;

	async function startQuiz() {
		const res = await fetch("http://localhost:8080/quiz/start", {
			method: "POST"
		});
		const data = await res.json();
		sessionId = data.session_id;
		question = data.question;
		explanation = "";
		status = "answering";
	}

	async function answer(option) {
		const res = await fetch("http://localhost:8080/quiz/answer", {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify({
				session_id: sessionId,
				question_id: question.id,
				topic_id: "articles",
				was_correct: option === getCorrectGuess(),
				difficulty: 2
			})
		});

		const data = await res.json();
		explanation = data.explanation;
		mastery = data.mastery;

		if (data.status === "finished") {
			status = "finished";
			question = null;
		} else {
			question = data.next_question;
		}
	}

	// TEMP: frontend guesses correctness (backend will validate later)
	function getCorrectGuess() {
		return question.options[0]; // just for MVP
	}
</script>

<h1>Chocolingo – Quiz MVP</h1>

{#if status === "idle"}
	<button on:click={startQuiz}>Start Quiz</button>
{/if}

{#if question}
	<h2>{question.prompt}</h2>

	{#each question.options as opt}
		<button on:click={() => answer(opt)}>
			{opt}
		</button>
	{/each}
{/if}

{#if explanation}
	<p><strong>Explanation:</strong> {explanation}</p>
{/if}

{#if mastery}
	<p>Mastery: {Math.round(mastery.Mastery)}%</p>
{/if}

{#if status === "finished"}
	<h2>Quiz finished 🎉</h2>
{/if}
