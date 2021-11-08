"use strict";
let ready = false;

const rooms = {};

let roomInView = "";

let nickname = "";

const errorModal = new bootstrap.Modal(document.getElementById("errorModal"));

const renderMessages = (messages) => {
  const messagesArea = document.getElementById("messages-area");
  if (messages && messages.length)
    for (let message of messages.reverse()) {
      if (message) {
        const messageElement = document.createElement("div");
        messageElement.className = "container-fluid my-4";
        messageElement.innerHTML = `<div class="container-fluid d-flex align-items-center"><img class="pfp message-pfp" src="./user.png" alt="user" />
        <div><p class="name">${
          message.displayName
        }<span class="time">${new Date(
          message.timestamp
        ).toLocaleTimeString()}</span></p><p>${
          message.message
        }</p></div></div>`;
        messagesArea.appendChild(messageElement);
      }
    }
};

const renderOnlineUsers = (users) => {
  const online = document.getElementById("online-users");
  online.innerHTML = "";
  for (let user of users) {
    const userElement = document.createElement("div");
    userElement.className =
      "d-flex mb-2 justify-content-center align-items-center";
    userElement.innerHTML = `<img class="pfp" src="./user.png" alt="user" /><div class="p-2 text-center text-white pointer">${user}</div>`;
    online.appendChild(userElement);
  }
};

const handleMessages = (message) => {
  rooms[message.room].messages.push({
    message: message.message,
    displayName: message.displayName,
  });
  if (message.room === roomInView) {
    renderMessages([message]);
  }
};

const renderRoom = (room) => {
  const roomNode = document.getElementById("room");
  roomNode.innerHTML = "";
  // add navbar
  const navbar = document.createElement("div");
  navbar.className = "navbar";
  navbar.innerHTML = `<div class="container-fluid"><span class="navbar-brand mb-0 h1">${room}</span></div>`;
  roomNode.appendChild(navbar);
  // add hr
  const hr = document.createElement("hr");
  hr.className = "m-0 mt-1";
  roomNode.appendChild(hr);
  const container = document.createElement("div");
  container.className = "row mw-100 h-100";
  container.innerHTML = `<div class="col-10 border-end d-flex flex-column" style="border-color: rgb(72, 72, 72) !important">
    <div id="messages-area" class="container-fluid mb-4 py-4" style="height: 92%"></div>
    <div class="container-fluid d-flex flex-column justify-content-center" style="height: 8%"><div id="message" class="input-group"></div></div>
    </div><div id="online" class="col-2 d-flex flex-column align-items-center p-2"><h4>Online</h4><hr class="my-2 w-100" /><div id="online-users"></div></div>`;
  roomNode.appendChild(container);
  const messageInputContainer = document.getElementById("message");
  const input = document.createElement("input");
  input.className = "form-control";
  input.type = "text";
  input.placeholder = "Type a message";
  input.addEventListener("keyup", (e) => {
    if (e.key === "Enter") {
      const message = input.value;
      input.value = "";
      if (message.length > 0) sendMessage(message, room);
    }
  });
  messageInputContainer.appendChild(input);
  const sendButton = document.createElement("button");
  sendButton.className = "btn btn-primary";
  sendButton.type = "button";
  sendButton.innerText = "Send";
  sendButton.addEventListener("click", () => {
    if (input.value.length > 0) sendMessage(message, room);
    input.value = "";
  });
  messageInputContainer.appendChild(sendButton);
  roomInView = room;
  renderMessages(...rooms[room].messages);
  renderOnlineUsers(rooms[room].users);
};

const renderRoomsList = () => {
  const roomList = document.getElementById("room-list");
  roomList.innerHTML = "";
  for (let room in rooms) {
    const roomElement = document.createElement("div");
    roomElement.className =
      "p-2 mb-4 text-center rounded-pill pe-auto bg-primary text-white pointer";
    roomElement.innerText = room;
    roomElement.addEventListener("click", () => renderRoom(room));
    roomList.appendChild(roomElement);
  }
};

const handleUserAdd = (message) => {
  if (!rooms[message.room]) {
    setTimeout(() => renderRoom(message.room), 600);
    rooms[message.room] = {
      users: [],
      messages: [],
      displayName: nickname,
    };
    renderRoomsList();
  }
  rooms[message.room].users.push(message.message);
  if (message.room === roomInView) renderOnlineUsers(rooms[message.room].users);
};

const handleHistory = (message) => {
  rooms[message.room].messages.push(
    message.history?.map((msg) => JSON.parse(msg))
  );
};

const handleUserRemove = (message) => {
  rooms[message.room].users = rooms[message.room].users.filter(
    (user) => user !== message.message
  );
  if (message.room === roomInView) renderOnlineUsers(rooms[message.room].users);
};

const handleErrors = (message) => {
  errorModal.toggle();
  document.querySelector("#error-message").innerText = message;
};

const processMessages = (message) => {
  if (!ready) return;
  switch (message.type) {
    case "UserAdd":
      handleUserAdd(message);
      break;
    case "UserRemove":
      handleUserRemove(message);
      break;
    case "History":
      handleHistory(message);
      break;
    case "MESSAGE":
      handleMessages(message);
      break;
    case "error":
      handleErrors(message.message);
      break;
  }
};

const socket = new WebSocket("ws://localhost:8080/ws");
socket.addEventListener("message", function (event) {
  if (event.data === "Connected!") {
    ready = true;
  } else {
    processMessages(JSON.parse(event.data));
  }
});

const createModal = new bootstrap.Modal(document.getElementById("createModal"));
const toggleCreateModal = () => createModal.toggle();
const createButton = document.querySelector("#createButton");
const roomNameInputCreate = document.querySelector("#roomNameToCreate");
const displayNameCreate = document.querySelector("#displayNameCreate");
const joinButton = document.querySelector("#joinButton");
const roomNameInputJoin = document.querySelector("#roomNameToJoin");
const displayNameJoin = document.querySelector("#displayNameJoin");
const joinModal = new bootstrap.Modal(document.getElementById("joinModal"));
const toggleJoinModal = () => joinModal.toggle();
joinButton.addEventListener("click", () => {
  if (!roomNameInputJoin.value || !displayNameJoin.value) return;
  console.log(roomNameInputCreate.value);
  socket.send(
    JSON.stringify({
      type: "JOIN",
      room: roomNameInputJoin.value,
      displayName: displayNameJoin.value,
      message: "Join room",
    })
  );
  nickname = displayNameJoin.value;
  toggleJoinModal();
});
createButton.addEventListener("click", () => {
  if (!roomNameInputCreate.value || !displayNameCreate.value) return;
  console.log(roomNameInputCreate.value);
  socket.send(
    JSON.stringify({
      type: "CREATE",
      room: roomNameInputCreate.value,
      displayName: displayNameCreate.value,
      message: "Create room",
    })
  );
  nickname = displayNameCreate.value;
  toggleCreateModal();
});

const sendMessage = (message, room) => {
  socket.send(
    JSON.stringify({
      type: "MESSAGE",
      timestamp: new Date().toISOString(),
      room: room,
      message: message,
      displayName: rooms[room].displayName,
    })
  );
};
