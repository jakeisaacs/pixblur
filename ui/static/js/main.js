// const socket = new WebSocket('ws://localhost:4000/ws');
            
// socket.onmessage = function(event) {
// 	console.log("event: ", event);
// 	const imagePath = event.data.split("|")[0]; 
// 	const timerValue = event.data.split("|")[1];  // The image file path sent from the server
// 	console.log('New image path:', imagePath);
// 	const imageElement = document.getElementById('live-image');
// 	const timerElement = document.getElementById('timer');

// 	timer.textContent = parseInt(timerValue, 10);
	
// 	// Update the image source with the new image
// 	imageElement.src = imagePath + '?' + new Date().getTime();  // Adding a timestamp to avoid caching issues
// };

// socket.onerror = function(error) {
// 	console.error('WebSocket error:', error);
// };

const wordLength = 6;

function handleKeyClick(e) {
	const keys = document.querySelectorAll('.key');
	const letterSlots = document.querySelectorAll('.letter-slot');
	const current_word = document.getElementById('current-word');
	const win_message = document.getElementById('win-message');

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
					win_message.classList.remove('hidden');
					socket.close()
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

