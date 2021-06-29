# Chaincode

## teamate.go structure
> managing users and projects

### data
User(
    id string
    avg float
    current_projects []string
    completed_projects []string
    <!-- proj_score_map map[int]int  use proj_id(int) as a key in the future -->
    proj_score_map map[string]int
)
Project(
    id string
    <!-- name string -->
    participants []string
    completed bool
)

### functions
registerUser 
    -> joinProject
    -> registerProject -> completeProject -> recordScore


## function details 
1. registerUser(id)
    -> initialize user
2. registerProject(proj_id,user_id)
    <!-- -> initialize project (proj_id : proj1, proj2...) -->

3. joinProject(proj_id,user_id)
    -> proj.participants.append(user_id)
    -> user.current_projects.append(proj_id)
4. completeProject(proj_id)
    -> proj.completed = true
    -> user.current_projects.remove(proj_id), user.completed_projects.append(proj_id)
5. recordScore(proj_id,user_id,score) 
    // check proj is completed
    -> user.proj_score_map[proj_id] = score
    <!-- -> user.proj_score_map[proj_id(int)] = score -->


## for the future update
* proj_name will be added
    > * registerProjcet use proj_name instead of proj_id : proj_id will automatically increase
    > * proj_score_map [string]int -> [int]int