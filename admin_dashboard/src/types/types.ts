export abstract class Question {
    id: number
    quizID: number
    type: string
    question: string
    next: number
    image: string
    answer: string[] = [""]


    constructor(
        id: number = 0,
        quizID: number = 0,
        type: string = "",
        question: string = "",
        next: number = 0,
        image: string = ""

    ) {
        this.id = id
        this.quizID = quizID
        this.type = type
        this.question = question
        this.next = next
        this.image = image
    }
}

export class QuestionText extends Question {
    answer: string[]

    constructor(
        id: number = 0,
        quizID: number = 0,
        type: string = "text",
        question: string = "",
        answer: string[] = [""],
        next: number = 0,
        image: string = ""
    ) {
        super(id, quizID, type, question, next, image)
        this.answer = answer
    }
}

export class QuestionSelect extends Question {
    answer: string[]
    options: string[]

    constructor(
        id: number = 0,
        quizID: number = 0,
        type: string = "select",
        question: string = "",
        next: number = 0,
        image: string = "",
        answer: string[] = [""],
        options: string[] = ["",""]
    ) {
        super(id, quizID, type, question, next, image)
        this.answer = answer
        this.options = options
    }
}

export class QuestionMultiSelect extends Question {
    answer: string[]
    options: string[]

    constructor(
        id: number = 0,
        quizID: number = 0,
        type: string = "multi_select",
        question: string = "",
        next: number = 0,
        image: string = "",
        answer: string[] = [],
        options: string[] = ["",""]
    ) {
        super(id, quizID, type, question, next, image)
        this.answer = answer
        this.options = options
    }
}

export class Quiz {
    id?: number
    type: string
    title: string
    description: string
    cover: string
    newAnswers: number = 0
    questions: Question[]

    constructor(
        id: number,
        type: string,
        title: string,
        description: string,
        cover: string,
        newAnswers: number = 0,
        questions: Question[] = []
    ) {
        this.id = id
        this.type = type
        this.title = title
        this.description = description
        this.cover = cover
        this.newAnswers = newAnswers
        this.questions = questions
    }
}
