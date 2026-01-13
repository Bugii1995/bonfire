<script>
	let sessionId = null;
	let question = null;
	let explanation = "";
	let mastery = null;
	let status = "idle"; // idle | answering | feedback | finished

	let selectedOption = null;
	let isCorrect = null;

	async function startQuiz() {
		const res = await fetch("http://localhost:8080/quiz/start", {
			method: "POST"
		});
		const data = await res.json();

		sessionId = data.session_id;
		question = data.question;

		explanation = "";
		selectedOption = null;
		isCorrect = null;
		status = "answering";
	}

	async function answer(option) {
		if (status !== "answering") return;

		// record what user clicked
		selectedOption = option;

		const res = await fetch("http://localhost:8080/quiz/answer", {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify({
				session_id: sessionId,
				question_id: question.id,
				topic_id: "articles",
				selected_option: option,
				difficulty: 2
			})
		});

		const data = await res.json();

		// backend decides truth
		isCorrect = data.is_correct;
		explanation = data.explanation;
		mastery = data.mastery;

		// ✅ ONLY now enter feedback state
		status = "feedback";

		// small pause so feedback is visible
		setTimeout(() => {
			if (data.status === "finished") {
				status = "finished";
				question = null;
			} else {
				question = data.next_question;
				selectedOption = null;
				isCorrect = null;
				explanation = "";
				status = "answering";
			}
		}, 900);
	}

	function buttonClass(opt) {
		if (status !== "feedback") return "";
		if (opt !== selectedOption) return "";
		if (isCorrect === null) return "";

		return isCorrect ? "correct" : "wrong";
	}
</script>

<style>
	body {
		font-family: system-ui, sans-serif;
		background: #fafafa;
	}

	.quiz {
		max-width: 480px;
		margin: 2rem auto;
		padding: 1.5rem;
		border: 1px solid #ddd;
		border-radius: 10px;
		background: white;
	}

	h1 {
		text-align: center;
	}

	button {
		display: block;
		width: 100%;
		margin: 0.5rem 0;
		padding: 0.7rem;
		font-size: 1rem;
		border-radius: 6px;
		border: 1px solid #ccc;
		cursor: pointer;
		transition: background 0.2s;
	}

	button:disabled {
		opacity: 0.7;
		cursor: not-allowed;
	}

	.correct {
		background: #c8f7c5;
		border-color: #3bb54a;
	}

	.wrong {
		background: #f7c5c5;
		border-color: #d33;
	}

	.explanation {
		margin-top: 1rem;
		padding: 0.75rem;
		background: #f5f5f5;
		border-left: 4px solid #888;
	}

	.mastery {
		margin-top: 0.5rem;
		font-size: 0.9rem;
		color: #444;
	}
</style>

<div class="quiz">
	<h1>Chocolingo</h1>

	{#if status === "idle"}
		<button on:click={startQuiz}>Start Quiz</button>
	{/if}

	{#if question}
		<h2>{question.prompt}</h2>

		{#each question.options as opt}
			<button
				class={buttonClass(opt)}
				disabled={status !== "answering"}
				on:click={() => answer(opt)}
			>
				{opt}
			</button>
		{/each}
	{/if}

	{#if explanation}
		<div class="explanation">
			<strong>{isCorrect ? "Correct!" : "Not quite."}</strong><br />
			{explanation}
		</div>
	{/if}

	{#if mastery}
		<div class="mastery">
			Mastery: {Math.round(mastery.Mastery)}%
		</div>
	{/if}

	{#if status === "finished"}
		<h2>Quiz finished 🎉</h2>
		<button on:click={startQuiz}>Restart</button>
	{/if}
</div>
