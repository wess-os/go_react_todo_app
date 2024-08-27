import { Flex, Spinner, Stack, Text } from "@chakra-ui/react";
import { useQuery } from "@tanstack/react-query";
import { BASE_URL } from "../App";
import TodoItem from "./TodoItem";

export type Todo = {
    _id: number;
    body: string;
    completed: boolean;
}

const TodoList = () => {
    const {data:todos, isLoading} = useQuery<Todo[]>({
        queryKey: ["todos"],

        queryFn: async () => {
            try {
                const res = await fetch(BASE_URL + "todos");
                const data = await res.json();
                
                if (!res.ok) {
                    throw new Error(data.error || "Something went wrong");
                }

                return data || [];
            } catch (error) {
                console.log(error);
            }
        }
    });

    return (
        <>
            {isLoading && (
                <Flex justifyContent={"center"} my={4}>
                    <Spinner size={"xl"} />
                </Flex>
            )}
            {!isLoading && todos?.length === 0 && (
                <Stack alignItems={"center"} gap='3' my={10}>
                    <Text fontSize={"xl"} textAlign={"center"} color={"gray.500"}>
                        Todas as tarefas foram concluídas!
                    </Text>
                </Stack>
            )}
            <Stack gap={3} my={10}>
                {todos?.map((todo) => (
                    <TodoItem key={todo._id} todo={todo as Todo} />
                ))}
            </Stack>
        </>
    );
}

export default TodoList;