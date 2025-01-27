function countDown() {
	let count = 3

	const gameOverlay = document.getElementById('gameOverlay');
	const countdownOverlay = document.getElementById('countdownOverlay');
	const countdown = document.getElementById('countdown');

	gameOverlay.classList.add('hidden');

	setTimeout(() => {
		countdown.textContent = count;
		countdownOverlay.classList.remove('hidden');

		const timer = setInterval(() => {
			count--;
			countdown.textContent = count;

			if (count <= 0) {
				console.log("starting game...")
				clearInterval(timer);
				countdownOverlay.classList.add('hidden');
				startGame();
			}
		}, 1000);
	}, 300);
}

function startGame() {
	const eventSource = new EventSource('http://localhost:4000/events');

	eventSource.onmessage = function(event) {
		const data = event.data.split('\n');  // The image file path sent from the server

		const timerValue = parseInt(data[0], 10);
		const imgSrc = data[1];

		const timerElement = document.getElementById('timer');
		const imageElement = document.getElementById('live-image');

		timerElement.textContent = timerValue;
		imageElement.src = imgSrc;  // Adding a timestamp to avoid caching issues
	};

	eventSource.onerror = function(error) {
		eventSource.close();
		console.error('Eventsource error:', error);
	};
}

function handleKeyClick(e) {
	const keys = document.querySelectorAll('.key');
	const letterSlots = document.querySelectorAll('.letter-slot');
	const current_word = document.getElementById('current-word');
	const winOverlay = document.getElementById('winOverlay');
	const wordLength = 6;

	let new_word = '';

	let currentSlotIndex = 0;
	letterSlots.forEach(slot => {
		new_word += slot.textContent;
		if (slot.textContent !== '') {
			currentSlotIndex++;
		} else {
			return;
		}
	});

	console.log("current: " + new_word);

	current_word.textContent = new_word;

	if (new_word.length === wordLength) {
		console.log("final: " + new_word);
		return;
	} 

	if (e.target.textContent === 'CLEAR') {
		letterSlots.forEach(slot => {
			slot.textContent = '';
		});
		currentSlotIndex = 0;
		return;
	}

	const key = e.target;
	const letter = key.textContent;
	const keyRect = key.getBoundingClientRect();
	const slotRect = letterSlots[currentSlotIndex].getBoundingClientRect();

	// Create floating letter
	const floatingLetter = document.createElement('div');
	floatingLetter.textContent = letter;
	floatingLetter.className = 'floating-letter';
	floatingLetter.style.left = `${keyRect.left + keyRect.width/2}px`;
	floatingLetter.style.top = `${keyRect.top + keyRect.height/2}px`;

	// Calculate translation
	const tx = (slotRect.left + slotRect.width/2) - (keyRect.left + keyRect.width/2);
	const ty = (slotRect.top + slotRect.height/2) - (keyRect.top + keyRect.height/2);
	floatingLetter.style.setProperty('--tx', `${tx}px`);
	floatingLetter.style.setProperty('--ty', `${ty}px`);

	document.body.appendChild(floatingLetter);

	// Add letter to slot
	setTimeout(() => {
		letterSlots[currentSlotIndex].textContent = letter;
		currentSlotIndex++;
		if (currentSlotIndex === wordLength) {
			new_word += letter;
			console.log("FINAL!!!: " + new_word);
			const data = {
				"word": new_word
			};
			const response = fetch('/check_word', {
				method: "POST",
				headers: {
					"Content-Type": "application/json"
				},
				body: JSON.stringify(data)
			});
			response.then(res => res.json()).then(data => {
				if (data["status"] === "success") {
					winOverlay.classList.remove('hidden');
					console.log("Correct word!");
				} else {
					console.log("Incorrect word!");
				}
			});
		} 
	}, 200);

	// Remove floating letter after animation
	setTimeout(() => {
		floatingLetter.remove();
	}, 100);

	// Remove key with animation
	key.classList.add('removing');
	setTimeout(() => {
		key.remove();
	}, 100);

}

