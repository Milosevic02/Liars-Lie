# Project Inspiration

Working on this project reminded me of the role of oracles in blockchain. Through the *Liars Lie* game task, I recognized a similarity in how data is collected and verified across a network. This challenge in the game felt like an analogy for filtering and assessing the reliability of data sources, which is a key aspect of using oracles.


## Start Functionality 

- **Processes for Agents**: Each agent is launched as a separate process (HTTP server) to simulate agent independence without shared memory, which would be the case if threads were used.

- **Structure**:
  - **`agent/agent.go`**: Defines the agent and its functionalities.
  - **`agent_starter/main.go`**: Launches the agent process with the appropriate parameters.
  - **`controllers/game.go`**: Contains the `StartGame` function, which is called when the `start` command is entered in **`cmd/main.go`**.

- **Logic**:
  - The `start` command initiates agents with parameters (`value`, `max-value`, `num-agents`, `liar-ratio`).
  - Before starting, `StopGame(false)` is called to ensure available ports.
  - **`StartGame`** launches agents as HTTP servers, assigning unique ports to each. Truthful agents return the correct value, while liars return a random value different from the correct one.
  - Agent ports are saved in the `agents.config` file.

- **Run command**: start --value v --max-value max --num-agents number --liar-ratio ratio

## Play Functionality

- **Communication with agents**: Each agent has a `/value` endpoint. The client sends a GET request to this endpoint to retrieve the value from each agent.  
  - The GET request is handled by the `handleGetValue` function in **`agent/server.go`**, which processes the request and returns the agent's value.

- **Value Processing**: The values returned by the agents are counted, and the "network value" is determined based on the most frequent value among the responses.
    - Although other algorithms were considered, such as calculating the average value or discarding extreme values, the counting approach for determining value repetitions proved to be the most efficient.

- **Client Logic**:
  - When the `play` command is received, the `Play` function in **`client/play.go`** is invoked.
  - This function sends HTTP GET requests to all agents, collects their responses, and processes them to find the most frequent value.
  - After processing the responses, the function then calculates and returns the network value.

- **Run command**: play

## Kill and Stop Functionality

- **Kill**:
  - Each agent has a `/kill` endpoint that handles a graceful shutdown.
  - When an agent is killed, the client will skip that agent during the `play` phase if no response is received. This ensures that the game continues without disruption if an agent unexpectedly shuts down.
  - The `kill` command is invoked when the client sends a request to the agent’s `/kill` endpoint, causing the agent to stop running gracefully.

- **Stop**:
  - The `stop` command stops all running agents by sending a `/kill` request to each of them. After that, it removes the `agents.config` file and deletes any log files, which are used for tracking events in the system.
  - If the `stop` command is called with the `true` flag, the entire program shuts down. This feature was implemented to allow the `stop` command to also be used before starting a new game (as described in the `start` section).

- **Run kill command**: kill --id id
- **Run stop command**: stop 

## Extend Functionality

### 1. **Value Issues (`value`)**
   - **Is it possible to change the initial value (`value`) after the game starts?**
   - To address this issue, several options were considered, and option 1 was selected and implemented. However, I would like to share a few alternatives:

   **Option 1: Adding new agents with the same value as at the start of the game**  
   - **Advantages**: This option feels natural because the new agents will be truthful and will have the appropriate number of liars. This approach maintains system consistency.  
   - **Disadvantages**: If the value (`value`) is changed, adding new agents with the same value can lead to undesirable effects, as all new agents will be set as liars, which can alter the game dynamics.

   **Option 2: Adding any agents according to the parameters**  
   - **Advantages**: This option fulfills the basic task, allowing any agents to be added according to the given parameters, without restrictions.  
   - **Disadvantages**: While it allows more flexibility, it could result in all newly added agents being liars, thus changing the fundamental setup of the game.

   **Option 3: Assigning a new value to all agents**  
   - **Advantages**: This option ensures that all truthful agents, regardless of the initial value, now have the same value, which enables greater consistency among the agents.  
   - **Disadvantages**: This option may result in agents losing the initial value that was set at the start, undermining the purpose of the initial value and the dynamics of the game.

---

### 2. **Liars Ratio Issues (`liar ratio`)**  
   - **How is the `liar ratio` handled when new agents are added?**
   - To address this problem, option 1 was chosen, but there are also alternatives that allow for different approaches:

   **Option 1: Calculating the `liar ratio` for newly added agents**  
   - **Advantages**: This approach is simple and direct, as the `liar ratio` is calculated only for the newly added agents. This is also in line with the specifications set in the task.  
   - **Disadvantages**: This approach can disrupt the accuracy of the initial `liar ratio` for the entire network of agents, as new agents may not necessarily fit into the existing ratio of liars and truthful agents.

   **Option 2: Equalizing the `liar ratio` for the entire network of agents**  
   - **Advantages**: This approach is more complex, but it offers greater accuracy as it attempts to maintain the balance of liars and truthful agents throughout the entire network, taking into account all agents, including new ones.  
   - **Disadvantages**: If the appropriate `liar ratio` is not provided when adding new agents (the same as at the start), the entire `liar ratio` may become unreliable, and the accuracy of the ratio between liars and truthful agents could be lost.

   - **Formula for calculating `liar ratio` for Option 2**:  
   The formula used to calculate whether an agent is a liar based on the `liar ratio` uses the number of agents in the network and the `liarRatio` value multiplied by 10. The formula is as follows:  
   `(len(agentProcesses)+i) % int(math.Round(liarRatio*10)) != 0`  
   - **Explanation**: The index of the newly added agent is divided by a value based on the `liarRatio` multiplied by 10. If the result of the division is zero, the agent is marked as a liar.

---

### 3. **Game Extension Logic (Extend)**  
   - When the appropriate command is recognized in `cmd/main.go`, the **Extend** function, located in **`controllers/game.go`**, is called.  
   - The logic is similar to the game startup, but with additional steps to update the agents' state and values to allow the addition of new agents, proper management of values, and maintenance of the correct liar ratio in the network.

- **Run command**: extend --value v --max-value max --num-agents number --liar-ratio ratio

## EXPERT Play

The **Expert Play** feature addresses the challenge of determining who the liars in the network are when a `liar ratio` is set. The key issue is that we don't initially know which agents are liars, and with the `liar ratio` given, it is difficult to choose the liars in a way that respects this ratio.

### Options Considered:

1. **Finding Liars through Agent Voting**
   - **Idea**: Each agent communicates with the remaining agents in the network, attempting to identify which ones it suspects to be liars. For example, if the `liar ratio` is set to 30% and there are 10 agents, the agent would return the 3 agents it suspects the most. Liars will vote for truth-tellers to mislead the system and try to eliminate truthful agents. Each agent reports the ones they suspect, and the client can then select the most suspicious agents (voting mechanism).  
   - **Idea for extension**: Add the possibility that if an agent is heavily voted against, it gets "killed" and removed from the network.  
   - **Advantages**: This approach introduces a solid mechanism for discovering liars and can improve the program’s functionality.  
   - **Disadvantages**: This doesn't fully solve the problem of how to select liars according to the specified `liar ratio` in the task. The idea could be adapted to select suspects for liars in a subgroup, but it does not align with the task’s requirements as it would introduce a voting mechanism to find suspects rather than respecting the `liar ratio` directly.

2. **Calling the Extend Function (Recreating Agents with Liar Ratio)**
Idea: Create new agents with the specified liar ratio. These new agents would communicate with the existing ones and return their values. When a new expert play command is issued, the previously created agents (the auxiliary ones created through the extend command) are killed, and new agents are created based on the updated `liar ratio`.
   - **Advantages**: Allows for the selection of liars according to the `liar ratio`.  
   - **Disadvantages**: Each time the `playexpert` function is called, new agents must be created, which could introduce unnecessary complexity. Additionally, these new agents would not have a predefined `value` and `max value` since these are not stored anywhere, and they would behave as clients gathering data, which complicates the liars’ behavior. The

 liars would also not have a predefined `value`, which contradicts the requirements of the task.  

3. **Randomly Selecting Agents**
   - **Advantages**: It would respect the structure of the task, except for the `liar ratio`.  
   - **Disadvantages**: Random selection does not guarantee that the `liar ratio` is respected, which is essential to the task. This is not an ideal solution since it doesn’t solve the core issue at hand.

4. **Maintaining a List of Known Liars**
   - **Idea**: Keep a list of agents identified as liars and select liars from that list.  
   - **Advantages**: It solves the problem of respecting the `liar ratio` since the liars are already known.  
   - **Disadvantages**: It’s unrealistic because the system won’t know who the liars are at the start, and it contradicts the distributed nature of the system. It also implies having a centralized system with all the data, which is not aligned with the task's goals.

### Most Interesting Option: **Option 1 – Agent Voting (Finding Liars)**
- **Reason for Interest**: I find this option the most intriguing because it introduces a dynamic mechanism for identifying liars based on communication and voting. Each agent communicates with others and votes for who they suspect to be a liar, which creates an interesting approach to managing liar detection in the network.  
- **Note**: The implementation for Option 1 can be found at the [**FindingLiars.md**](FindingLiars.md)


### Chosen Option: **Option 3 – Random Agent Selection**
- **Reason for Selection**: This option aligns most closely with the task requirements. Although it does not fully respect the `liar ratio`, it is the most straightforward approach for integrating with the system and adhering to the task's guidelines.  

### Implementation Logic:
- When the **`playexpert`** command is invoked from `cmd/main.go`, it calls the function defined in `client/play_expert.go`.  
- Each agent has an endpoint for this operation and receives a list of port numbers for agents outside of the selected subset. The agents then interact with the others, and the liars will only communicate their values without further processing.  
- The truthful agents will calculate the correct value and return it to the client.

- **Run command**: play-expert --num-agents number --liar-ratio ratio

# Tests

The tests for this project are located in the `tests` folder. There is also a shell script called `start_all.sh` which runs all the tests automatically.
