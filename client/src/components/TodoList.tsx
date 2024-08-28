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
    const token = localStorage.getItem('token');

    const {data:todos, isLoading} = useQuery<Todo[]>({
        queryKey: ["todos"],

        queryFn: async () => {
            try {
                const headers: HeadersInit = {
                    "Content-Type": "application/json",
                };

                if (token) {
                    headers["Authorization"] = token;
                }

                const res = await fetch(BASE_URL + "todos", {
                    method: 'GET',
                    headers: headers,
                });
                const data = await res.json();
                
                if (!res.ok) {
                    throw new Error(data.error || "Something went wrong");
                }

                return data || [];
            } catch (error) {
                window.location.href = '/login';
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